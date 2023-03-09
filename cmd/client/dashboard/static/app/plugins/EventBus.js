app.use({
    install: (app, _) => {
        let listeners = [], 
            global = app.config.globalProperties

        global.$bus = {
            trigger: (eventName, data) => {
                if (listeners[eventName] !== void 0) {
                    listeners[eventName].forEach(fn => fn(data))
                }
            },

            once: (eventName, fn) => {
                let onceFn = (record) => {
                    const index = listeners[eventName].indexOf(onceFn)
                    listeners[eventName].splice(index, 1)
                    fn(record)
                }
                global.$bus.on(eventName, onceFn)
            },

            on: (eventName, fn) => {
                if (listeners[eventName] === void 0) {
                    listeners[eventName] = []
                }

                listeners[eventName].push(fn)
            }
        }
    }
})