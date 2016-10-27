package event

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"sync"
)

// Default number of maximum listeners for an event.
const DefaultMaxListeners = 10

// Error presented when an invalid argument is provided as a listener function
var ErrNoneFunction = errors.New("Kind of Value for listener is not Func.")

// RecoveryListener ...
type RecoveryListener func(event interface{}, args []interface{}, listener interface{}, err error)
// If there is no listener for event
type DeadEventListener func(event interface{}, args []interface{})

type Emitter interface {
	Once(event, listener interface{}) Emitter
	On(event, listener interface{}) Emitter
	Off(event, listener interface{}) Emitter
	Emit(event interface{}, arguments ...interface{}) Emitter
	//GetListenerCount(event interface{}) (count int)
	//SetMaxListeners(max int) Emitter
	//RecoverWith(listener RecoveryListener) Emitter
	EmitSync(event interface{}, arguments ...interface{}) Emitter
}
// Emitter ...
type StdEmitter struct {
	Emitter
	// Mutex to prevent race conditions within the Emitter.
	*sync.Mutex
	// Map of event to a slice of listener function's reflect Values.
	events            map[interface{}][]reflect.Value
	// Optional RecoveryListener to call when a panic occurs.
	recoverer         RecoveryListener
	deadEventListener DeadEventListener
	// Maximum listeners for debugging potential memory leaks.
	maxListeners      int
	// Map used to remove Listeners wrapped in a Once func
	onces             map[reflect.Value]reflect.Value
}

// AddListener appends the listener argument to the event arguments slice
// in the Emitter's events map. If the number of listeners for an event
// is greater than the Emitter's maximum listeners then a warning is printed.
// If the relect Value of the listener does not have a Kind of Func then
// AddListener panics. If a RecoveryListener has been set then it is called
// recovering from the panic.
func (emitter *StdEmitter) AddListener(event, listener interface{}) *StdEmitter {
	emitter.Lock()
	defer emitter.Unlock()

	fn := reflect.ValueOf(listener)

	if reflect.Func != fn.Kind() {
		handlerRecover(emitter.recoverer, event, listener, ErrNoneFunction, nil)
	}

	if emitter.maxListeners != -1 && emitter.maxListeners < len(emitter.events[event]) + 1 {
		fmt.Fprintf(os.Stdout, "Warning: event `%v` has exceeded the maximum " +
			"number of listeners of %d.\n", event, emitter.maxListeners)
	}

	emitter.events[event] = append(emitter.events[event], fn)

	return emitter
}

// On is an alias for AddListener.
func (emitter *StdEmitter) On(event, listener interface{}) Emitter {
	return emitter.AddListener(event, listener)
}

// RemoveListener removes the listener argument from the event arguments slice
// in the Emitter's events map.  If the reflect Value of the listener does not
// have a Kind of Func then RemoveListener panics. If a RecoveryListener has
// been set then it is called after recovering from the panic.
func (emitter *StdEmitter) RemoveListener(event, listener interface{}) *StdEmitter {
	emitter.Lock()
	defer emitter.Unlock()

	fn := reflect.ValueOf(listener)

	if reflect.Func != fn.Kind() {
		handlerRecover(emitter.recoverer, event, listener, ErrNoneFunction, nil)
	}

	if events, ok := emitter.events[event]; ok {
		if _, ok = emitter.onces[fn]; ok {
			fn = emitter.onces[fn]
		}

		for i, listener := range events {
			if fn == listener {
				// Do not break here to ensure the listener has not been
				// added more than once.
				emitter.events[event] = append(emitter.events[event][:i], emitter.events[event][i + 1:]...)
			}
		}
	}

	return emitter
}
// General recover process
func handlerRecover(recoverer RecoveryListener, event interface{}, listener interface{}, r interface{}, args[]interface{}) {
	if r == nil {
		// Nothing to recover
		return
	}
	var err error
	if _, ok := r.(error); ok {
		err = r.(error)
	} else {
		err = fmt.Errorf("%v", r)
	}
	if recoverer == nil {
		panic(err)
	}
	recoverer(event, args, listener, err)
}

func (emitter *StdEmitter)deadEvent(event interface{}, args[] interface{}) {
	listener := emitter.deadEventListener
	if listener == nil {
		return
	}
	recoverer := emitter.recoverer
	if recoverer != nil {
		defer func() {
			if r := recover(); r != nil {
				handlerRecover(recoverer, emitter, listener, r, args)
			}
		}()
	}
	listener(event, args)
}

// Off is an alias for RemoveListener.
func (emitter *StdEmitter) Off(event, listener interface{}) Emitter {
	return emitter.RemoveListener(event, listener)
}

// Once generates a new function which invokes the supplied listener
// only once before removing itself from the event's listener slice
// in the Emitter's events map. If the reflect Value of the listener
// does not have a Kind of Func then Once panics. If a RecoveryListener
// has been set then it is called after recovering from the panic.
func (emitter *StdEmitter) Once(event, listener interface{}) Emitter {
	fn := reflect.ValueOf(listener)
	if reflect.Func != fn.Kind() {
		handlerRecover(emitter.recoverer, event, listener, ErrNoneFunction, nil)
	}

	var run func(...interface{})

	run = func(arguments ...interface{}) {
		defer emitter.RemoveListener(event, run)

		var values []reflect.Value

		for i := 0; i < len(arguments); i++ {
			values = append(values, reflect.ValueOf(arguments[i]))
		}

		fn.Call(values)
	}

	// Lock before changing onces
	emitter.Lock()
	emitter.onces[fn] = reflect.ValueOf(run)
	emitter.Unlock()

	emitter.AddListener(event, run)
	return emitter
}

// Emit attempts to use the reflect package to Call each listener stored
// in the Emitter's events map with the supplied arguments. Each listener
// is called within its own go routine. The reflect package will panic if
// the arguments supplied do not align the parameters of a listener function.
// If a RecoveryListener has been set then it is called after recovering from
// the panic.
func (emitter *StdEmitter) Emit(event interface{}, arguments ...interface{}) Emitter {
	var (
		listeners []reflect.Value
		ok bool
	)

	// Lock the mutex when reading from the Emitter's
	// events map.
	emitter.Lock()

	if listeners, ok = emitter.events[event]; !ok || len(listeners) == 0 {
		// If the Emitter does not include the event in its
		// event map, it has no listeners to Call yet.
		emitter.Unlock()
		// There is no listener for this event
		emitter.deadEvent(event, arguments)
		return emitter
	}

	// Unlock the mutex immediately following the read
	// instead of deferring so that listeners registered
	// with Once can acquire the mutex for removal.
	emitter.Unlock()

	var wg sync.WaitGroup

	wg.Add(len(listeners))
	recoverer := emitter.recoverer
	for _, fn := range listeners {
		go func(fn reflect.Value) {
			// Recover from potential panics, supplying them to a
			// RecoveryListener if one has been set, else allowing
			// the panic to occur.
			if nil != recoverer {
				defer func() {
					if r := recover(); nil != r {
						handlerRecover(recoverer, event, fn.Interface(), r, arguments)
					}
				}()
			}

			var values []reflect.Value

			for i := 0; i < len(arguments); i++ {
				if arguments[i] == nil {
					values = append(values, reflect.New(fn.Type().In(i)).Elem())
				} else {
					values = append(values, reflect.ValueOf(arguments[i]))
				}
			}

			defer wg.Done()

			fn.Call(values)
		}(fn)
	}

	wg.Wait()
	return emitter
}

// EmitSync attempts to use the reflect package to Call each listener stored
// in the Emitter's events map with the supplied arguments. Each listener
// is called synchronously. The reflect package will panic if
// the arguments supplied do not align the parameters of a listener function.
// If a RecoveryListener has been set then it is called after recovering from
// the panic.
func (emitter *StdEmitter) EmitSync(event interface{}, arguments ...interface{}) Emitter {
	var (
		listeners []reflect.Value
		ok bool
	)

	// Lock the mutex when reading from the Emitter's
	// events map.
	emitter.Lock()

	if listeners, ok = emitter.events[event]; !ok || len(listeners) == 0 {
		// If the Emitter does not include the event in its
		// event map, it has no listeners to Call yet.
		emitter.Unlock()
		// There is no listener for this event
		emitter.deadEvent(event, arguments)
		return emitter
	}

	// Unlock the mutex immediately following the read
	// instead of deferring so that listeners registered
	// with Once can aquire the mutex for removal.
	emitter.Unlock()

	recoverer := emitter.recoverer

	for _, fn := range listeners {
		var values []reflect.Value

		for i := 0; i < len(arguments); i++ {
			if arguments[i] == nil {
				values = append(values, reflect.New(fn.Type().In(i)).Elem())
			} else {
				values = append(values, reflect.ValueOf(arguments[i]))
			}
		}
		if recoverer == nil {
			fn.Call(values)
		} else {
			// Recover from potential panics, supplying them to a
			// RecoveryListener if one has been set, else allowing
			// the panic to occur.
			func() {
				defer func() {
					if r := recover(); nil != r {
						handlerRecover(recoverer, event, fn.Interface(), r, arguments)
					}
				}()
				fn.Call(values)
			}()
		}
	}

	return emitter
}

// RecoverWith sets the listener to call when a panic occurs, recovering from
// panics and attempting to keep the application from crashing.
func (emitter *StdEmitter) RecoverWith(listener RecoveryListener) *StdEmitter {
	emitter.recoverer = listener
	return emitter
}

// SetDeadEventListener sets the listener to call when there is no listener for an event
func (emitter *StdEmitter) SetDeadEventListener(listener DeadEventListener) *StdEmitter {
	emitter.deadEventListener = listener
	return emitter
}

// SetMaxListeners sets the maximum number of listeners per
// event for the Emitter. If -1 is passed as the maximum,
// all events may have unlimited listeners. By default, each
// event can have a maximum number of 10 listeners which is
// useful for finding memory leaks.
func (emitter *StdEmitter) SetMaxListeners(max int) *StdEmitter {
	emitter.Lock()
	defer emitter.Unlock()

	emitter.maxListeners = max
	return emitter
}

// GetListenerCount gets count of listeners for a given event.
func (emitter *StdEmitter) GetListenerCount(event interface{}) (count int) {
	emitter.Lock()
	if listeners, ok := emitter.events[event]; ok {
		count = len(listeners)
	}
	emitter.Unlock()
	return
}

// NewEmitter returns a new Emitter object, defaulting the
// number of maximum listeners per event to the DefaultMaxListeners
// constant and initializing its events map.
func NewEmitter() (emitter *StdEmitter) {
	emitter = new(StdEmitter)
	emitter.Mutex = new(sync.Mutex)
	emitter.events = make(map[interface{}][]reflect.Value)
	emitter.maxListeners = DefaultMaxListeners
	emitter.onces = make(map[reflect.Value]reflect.Value)
	return
}

// Default emitter
var DefaultEmitter = Emitter(NewEmitter().SetMaxListeners(30))

func Once(event, listener interface{}) Emitter {
	return DefaultEmitter.Once(event, listener)
}
func On(event, listener interface{}) Emitter {
	return DefaultEmitter.On(event, listener)
}
func Off(event, listener interface{}) Emitter {
	return DefaultEmitter.Off(event, listener)
}
func Emit(event interface{}, arguments ...interface{}) Emitter {
	return DefaultEmitter.Emit(event, arguments)
}
func EmitSync(event interface{}, arguments ...interface{}) Emitter {
	return DefaultEmitter.EmitSync(event, arguments)
}
