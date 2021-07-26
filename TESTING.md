# Testing

Code in this repository will run against operational bare metal infrastructure, and must be held to a higher standard than code that can be vetted as it progresses through similar environments. Consumers of this product will not verify compatibiity with their potential edge or corner cases: regressions and behavioral changes must be avoided. As this code will run, but not be tested, on actual infrastructure, thorough testing is paramount. The testing guidelines below outline our commitment to quality, consistency, and reliability. Pull requests that do not satisfy these guidelines are works in progress and will be handled as such.

## General Guidelines

- Use `testing.T` for testing code in lieu of adding extra dependencies
- [Prefer table driven tests](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)
- **Document behavior through testing**. The article linked above demonstrates this concept.
- Test all error conditions. Your tests should convey when, why, and how errors occur.
- Run `make .git/hooks/pre-commit` to install a pre-commit hook to run `make check`

When in doubt, use existing code as an example.

## Types of Tests

### Unit Testing

All functions and methods must have basic unit testing. For simple functions such as

```golang
// pkg/maas/machine.go

// SystemID returns the machine's systemID.
func (m *MachineManager) SystemID() string {
    return m.Current().SystemID
}
```

even the simplest test protects against regressions due to changes in `MachineManager.Current()` and `Machine.SystemID`. An acceptable test for this method would:

- Create a `MachineManager`
- Insert multiple `Machine`s with different `SystemID`s, storing the last (ie `Current`) one
- Verify the returned value matches the stored value **without** calling `Current()`

Functions that manipulate variables such as their parameters, their receiver, or attributes of their receiver must be tested against invalid values to prevent a [cadre of unexpected behaviors](https://github.com/dvyukov/go-fuzz#trophies). Tests also act as documentation for nontrivial usage:

```golang
func StripFirstThreeCharacters(s string) string {
    return s[3:]
}
StripFirstThreeCharacters("it") // Oops
StripFirstThreeCharacter("✗-men") // What does this return?

type MyTypes []MyType
func (m MyTypes) GetElementAtIndex(idx int) MyType {
    return m[idx]
}
var mt MyTypes
mt.GetElementAtIndex(-1) // Yikes
mt.GetElementAtIndex(1<<63) // And the other side
```

A simple function is a function without control flow (`if`, `for`, `select`, `switch`, `defer`). Functions with control flow must test all code paths in order to:

- Demonstrate flow to each code path
- Verify sound logic (difficult testing indicates unsuitable logic/structure)
- Protect against regressions by "documenting" expected behavior

Functions using `recover()` must demonstrate a `panic()`, and those that `defer` anonymous functions must demonstrate the logic in those functions; if they access variables from an outside scope, they are subject to the requirements described above.

Functions that return errors must demonstrate how to identify the error, how to cause the error, and whether the error is fatal:

```golang
func Foo(i, j int) error {
    if i > j {
        // Can execution continue without more j?
        return errors.New("Not enough j")
    }
    if i == j {
        // - Does this return a SameValueError?
        // - What is the message of this error?
        // - Is there any simple introspection that can be done on this error
        //     (eg err.TheValue == i == j, etc)?
        return NewSameValueError(i, j)
    }
    // Why would this return an error?
    if err := doStuffWith(i, j); err != nil {
        return err // What is err?
    }
}
```

Functions with goroutines, IPC such as channels, and process control such as mutexes and wait groups must test that functionality. These tests should demonstrate these tools were used safely:

- Goroutines must be stoppable (eg with a `Context` or `<-done` on a select)
- IPC must be protected against issues such as infinite blocking, sending on closed channels, and deadlocks
- Process control must be safeguarded against similar issues, such as a missing `wg.Done(1)` or calling `myMutex.Lock()` before calling a function that locks the same mutex

Tests that do not clearly demonstrate the behavior of the functions they describe will be considered incomplete.
