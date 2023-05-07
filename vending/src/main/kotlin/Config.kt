package main

data class Config(
    val debug: Boolean,
    val emulationMode: Boolean
)

fun parseConfig(args: Array<String>) =
    Config(
        debug = args.find { it == "--debug" } != null,
        emulationMode = args.find { it == "--emulation" } != null
    )

