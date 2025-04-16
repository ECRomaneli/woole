app.provide('$timer', {
    debounce(fn, delay) {
        let timeout
        return function (...args) {
            clearTimeout(timeout)
            timeout = setTimeout(() => {
                fn.apply(this, args)
            }, delay)
        }
    }
})