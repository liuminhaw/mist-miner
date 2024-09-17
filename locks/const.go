package locks

const (
	LINUX_DIR_PATH = "/var/lock/mist-miner"

	OBJECTS_LOCKFILE         = "mm-objects.lock"
	HISTORY_LOCKFILE         = "mm-history-logger.lock"
	HISTORY_POINTER_LOCKFILE = "mm-history-pointer.lock"
	REF_MARK_LOCKFILE        = "mm-refmark.lock"

	lock_retry_max_interval = 500
	lock_retry_max_count    = 5
)
