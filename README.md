
# Session Logger

This package adds a "session logger", a log system that allows you to create loggers with three log levels, a
prefix, and a random unique ID that will be used to prefix every message. This makes it easy to tell where the
messages are coming from, if a request generates multiple log messages, etc.

All the session loggers for a given program share a log file, and also log to stdout and stderr.

This is not intended to be the One True Logging Solution™ rather I made this to simplify basic logging in
simple server applications, specifically in REST endpoints for low traffic server apps, chat bots, and other
endpoint or callback based microservices.

Under the covers, logging is done by the standard library log package. No attempt is made to control log file
size outside of creating a new one every time the program is started. You can disable logging to a file if the
program will be run in a container or other system that manages turning standard out/err into log files for you.
