# Melee

Integration and load testing with simulation of the Clients Flow.

## Quick start

## Structure of tests

1. Basic element of tests is step. Step is a structure:

    ```go
    type Step struct {
        Name string
        Func func(f *Flow)
    }
    ```

    You can create new step with func `NewStep(stepName, stepFunc)`,
    where `stepName` is a string and `stepFunc` is `func(f *Flow)`,
    where `f` is a flow to which this step belongs.

    ```go
    func StepRegistration(f *Flow) {
        c := NewClient(f.Client)
        assert.Nil(f, c.Provider.NewSession())

        u, err := c.CreateRandomUser()
        assert.Nil(f, err)
        f.SetValue(user, u)
    }
    ```

2. Second element is flow:

    ```go
    type Flow struct {
        Name              string
        ContextKV         map[CtxKey]interface{}
        ctxRWMutex        sync.RWMutex
        Client            SessionProvider
        ETHClientProvider EthClientProvider
        steps             []Step

        *testing.T
    }
    ```

    Example:

    ```go
    func Login() *flow.Flow {
        provider := session.NewProvider()
        return flow.NewFlow("User registration and login.").SetSessionProvider(&provider).AddSteps(
            flow.NewStep("Registration", StepRegistration),
            flow.NewStep("Login", StepLogin),
        )
    }
    ```

## Highload testing

Coming soon
