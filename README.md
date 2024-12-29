# DGCommander

A [bwmarrin/discordgo](https://github.com/bwmarrin/discordgo) library for handling application commands. It routes commands to handlers and provides utilities to answer to these commands.

## How it works

Commands are defined using the builders,
these builders allow to know everything that a command can have without
needing to check the discord developer docs.
Then, the builders are given to a `DGCommander` instance that will register the commands
and handle argument parsing and routing of the command to the given handler.

If the arguments are invalid, for example a required argument was not provided,
then an error is returned to the user and the handler is never invoked.
This way we can trust that if we specify an argument to have a shape,
by the time we get to the handler, an argument with that shape will be present

```go
// Mashup from the /examples folder
commander := dgc.New(slog.Default(), session, dgc.DefaultTimeProvider{})

commander.AddCommand(
    dgc.NewMessageCommand(). // Commands on messages
        Name("Resend message").
        AllowEverywhere(true).
        Handler(func(ctx *dgc.MessageExecutionContext) error {
                defer ctx.Finish() // Stops the internal timer that updates the ctx to know if an interaction is still valid
                return ctx.RespondWithMessage(&discordgo.InteractionResponseData{
                    Content: ctx.Message.Content,
                    Embeds:  ctx.Message.Embeds,
                })
            }),
)
commander.AddCommand(
    dgc.NewUserCommand(). // Commands on users
        Name("User information").
        AllowEverywhere(true).
        Handler(func(ctx *dgc.UserExecutionContext) error {
                defer ctx.Finish()
                return ctx.RespondWithMessage(&discordgo.InteractionResponseData{
                    Content: "User information",
                    Embeds: []*discordgo.MessageEmbed{{
                        Title: ctx.User.Username,
                        Color: ctx.User.AccentColor,
                        Image: &discordgo.MessageEmbedImage{
                            URL: ctx.User.AvatarURL(""),
                        },
                        Fields: []*discordgo.MessageEmbedField{{
                            Name:  "Bot?",
                            Value: fmt.Sprintf("%t", ctx.User.Bot),
                        }},
                    }},
                })
            }),
)
commander.AddCommand(
    dgc.NewMultiSlashCommandBuilder(). // Chat/Slash commands
        Name("calculate").
        Description("Collection of simple operations").
        AllowEverywhere(true).
        AddSubCommand(dgc.NewSubCommand().
            Name("sum").
            Description("Calculates the sum of 2 numbers `a+b`").
            AddArguments(
                dgc.NewNumberArgument().Name("a").Description("First value of the sum").Required(true),
                dgc.NewNumberArgument().Name("b").Description("Seccond value of the sum").Required(true),
            ).
            Handler(func(ctx *dgc.SlashExecutionContext) error {
                    defer ctx.Finish()
                    a := ctx.GetRequiredNumber("a")
                    b := ctx.GetRequiredNumber("b")
                    return ctx.RespondWithMessage(&discordgo.InteractionResponseData{
                        Content: fmt.Sprintf("The result of `%.2f + %.2f` is `%.2f`", a, b, a+b),
                    })
                }),
        ).
        AddSubCommandGroup(dgc.NewSubCommandGroup().
            Name("trigonometry").
            Description("Trigonometry related functions").
            AddSubCommand(dgc.NewSubCommand().
                Name("sin").
                Description("Calculates the sin of the given angle").
                AddArguments(
                    dgc.NewNumberArgument().Name("angle").Description("The angle to calculate the sin").Required(true),
                    dgc.NewNumberChoicesArgument().Name("degree").Description("Degree of type to calculate").Required(false).
                        AddChoice("degrees", math.Pi/180.01).
                        AddChoice("radians", 1),
                ).
                Handler(func(ctx *dgc.SlashExecutionContext) error {
                        defer ctx.Finish()
                        angle := ctx.GetRequiredNumber("angle")
                        degree := ctx.GetNumberOr("degree", 1) // default is radians
                        angle *= degree
                        return ctx.RespondWithMessage(&discordgo.InteractionResponseData{
                            Content: fmt.Sprintf("The `sin(%.2f)` is `%.2f`", angle, math.Sin(angle)),
                        })
                    }),
            ),
        ),
)
```

## Why?

I have several discord bots in mind that I want to make. Ideally, I wish to use Golang to make them.
When searching for libraries to make them, I didn't find them, didn't like what I found or didn't solve the problem I was trying to fix, so I started making them myself.

This is one of several libraries I'm making for bot development in golang using discordgo, but I'm working on more like:

- [go-n-i18n](https://github.com/MrNemo64/go-n-i18n) for internationalization. While oriented to bot development, it's done to be usable in any golang application.
- [DGCommander](https://github.com/MrNemo64/dgcommander) for command handling.
- Widget library to make widgets using messages, embeds, and message components.
- Pagination library to show data.

I develop them as I need them for my projects so while I wish to make them all, I won't start development untill I need them.

Is it the best library for commands? Probably not, but it has/will have all that I want a command library to have:

1. Allows me to define a command without having to open the developer docs to know that command type had what.
1. Allows me to define a shape for a command and a handler, not having to worry about parsing arguments or whatever.

   - This includes autocompleting arguments, allowing me to just specify a handler that will return the autocomplete options for a given argument (see the tasks example).

1. Allows me to "personalize" commands, allowing me to create and reuse custom argument types like an argument for a duration (see `dgc/extras/duration_argument.go`, used in the tasks example).
   - Will also add middleware to commands so extra steps can be made even before the handler of a command is invoked.

## Features

This library is still very early in development and can change at any time as I realise that changes need to be made when using it.

- Support for **message, user and slash application commands**. Most libraries I've found just use text messages for commands and use a prefix to differenciate between a normal message and a command message. DGCommander does not support message based commands, only application commands as defined by the [Discord docs](https://discord.com/developers/docs/interactions/application-commands#application-commands).

- Use of **builder pattern to allow to create commands** without needing to read the Discord docs to know what each kind of command needs (for example, slash commands have a description but user and message commands don't).

- **Argument verification and parsing** for slash commands.

- Allows to **define custom argument types** (see `dgc/extras/duration_argument.go` for an example).

- Methods to easily interact with a command, like getting the arguments of a slash command, responding to a command, obtaining the time at which the interaction token expires...

- Use of the `context` package to be able to stop a handler if the valid time of an interaction expires.

- (Not done yet) Allow the use of **middleware** to transform/stop a command before it reaches the handler.
