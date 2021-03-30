# Powerplan2go

A simple windows 10 tray icon application for managing the current powerplan used.

Why? Because on a desktop computer I can't find an easy way to show the windows builtin
menu and I'm switching powerplans regularly.

## How does it work?

In the background it's just calling `powercfg` with the corresponding CLI-flags.
