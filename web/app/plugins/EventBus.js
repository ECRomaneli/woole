app.use({
    install: (app) => {
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
                return global.$bus.on(eventName, onceFn)
            },

            on: (eventName, fn) => {
                if (listeners[eventName] === void 0) {
                    listeners[eventName] = []
                }

                listeners[eventName].push(fn)
                return fn
            },

            off: (eventName, fn) => {
                const list = listeners[eventName]
                if (list === void 0) {
                    console.warn("There is no listeners for event " + eventName)
                    return
                }

                if (fn === void 0) {
                    listeners[eventName] = []
                    return
                }

                const indexToRemove = list.indexOf(fn)
                if (indexToRemove === -1) {
                    console.warn("Listener not found, ignoring...")
                    return
                }

                list.splice(indexToRemove, 1)
            }
        }
    }
})