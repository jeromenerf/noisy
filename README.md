# Noisy

> Makes some annoying clicky noise on your keyboard

`noisy` takes inspiration from [atom's mechanical keyboard package](https://atom.io/packages/mechanical-keyboard) and [dddpaul's go-evhandler](https://github.com/dddpaul/go-evhandler).

It plays mechanical keyboard sounds (from [http://www.freesfx.co.uk/](http://www.freesfx.co.uk/)) as you type. There are actually three samples.

Install with `go get github.com/jeromenerf/noisy`. Dependencies and samples are vendored in. It should be possible to configure it using spf13's viper.

As is, it requires `mpg123` to play mp3 files.
