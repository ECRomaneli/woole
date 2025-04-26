app.provide('$timer', {
    debounce(fn, delay) {
        let timeout
        return function (...args) {
            clearTimeout(timeout)
            timeout = setTimeout(() => { fn.apply(this, args) }, delay)
        }
    },

    debounceWithThreshold(fn, delay) {
        let lastArgs = null
        let isScheduled = false

        const execute = () => {
            if (lastArgs) {
                lastArgs = null
                fn.apply(this, lastArgs)
            }
        }

        return function (...args) {
            lastArgs = args

            if (!isScheduled) {
                isScheduled = true
                timeout = setTimeout(() => { isScheduled = false; execute() }, delay)
            }
        }
    }
})