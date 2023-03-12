app.use({
    install: (app, _) => {
        const STRING = 'string', NUMBER = 'number', BOOLEAN = 'boolean'

        function clone(obj) {
            const typeOf = typeof obj

            if (typeOf === STRING || typeOf === NUMBER || typeOf === BOOLEAN || obj === null || obj === void 0) {
                return obj 
            } 
            
            if (Array.isArray(obj)) {
                const listClone = []
                for (const item of obj) { listClone.push(clone(item)) }
                return listClone
            }

            const objClone = {}
            for (const key of Object.keys(obj)) { objClone[key] = clone(obj[key]) }
            return objClone
        }

        app.provide('$clone', clone)
    }
})