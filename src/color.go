package main

type Color string

const (
    Red     Color = "\033[31m"
    Green   Color = "\033[32m"
    Yellow  Color = "\033[33m"
    Blue    Color = "\033[34m"
    Magenta Color = "\033[35m"
    Cyan    Color = "\033[36m"
    White   Color = "\033[37m"
    Bold    Color = "\033[1m"
    Reset   Color = "\033[0m"
)

func color(text string, c Color) string {
    return string(c) + text + string(Reset)
}
