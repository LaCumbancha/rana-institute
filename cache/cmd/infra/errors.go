package infra

type CacheErrors string
const NOT_FOUND CacheErrors = "NOT_FOUND"
const EXPIRED CacheErrors = "EXPIRED"
const OK CacheErrors = ""
