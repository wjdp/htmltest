package htmltest

type Options struct {
  EnforceHTTPS bool
}

func NewOptions() Options {
  // Specify defaults here
  options := Options{
    EnforceHTTPS: false,
  }
  return options
}
