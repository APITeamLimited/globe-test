package logrus

// A hook to be fired when logging on the logging levels returned from
// `Levels()` on your implementation of the interface. Note that this is not
// fired in a goroutine or a channel with workers, you should handle such
// functionality yourself if your call is non-blocking and you don't wish for
// the logging calls for levels returned from `Levels()` to block.
type Hook interface ***REMOVED***
	Levels() []Level
	Fire(*Entry) error
***REMOVED***

// Internal type for storing the hooks on a logger instance.
type LevelHooks map[Level][]Hook

// Add a hook to an instance of logger. This is called with
// `log.Hooks.Add(new(MyHook))` where `MyHook` implements the `Hook` interface.
func (hooks LevelHooks) Add(hook Hook) ***REMOVED***
	for _, level := range hook.Levels() ***REMOVED***
		hooks[level] = append(hooks[level], hook)
	***REMOVED***
***REMOVED***

// Fire all the hooks for the passed level. Used by `entry.log` to fire
// appropriate hooks for a log entry.
func (hooks LevelHooks) Fire(level Level, entry *Entry) error ***REMOVED***
	for _, hook := range hooks[level] ***REMOVED***
		if err := hook.Fire(entry); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***
