app.provide('$map', {    
    toKeyValuePairs: (map) => {
        if (map === void 0) { return [] }
        
        const keyValuePairs = []
        Object.keys(map).forEach(key => {
            // Skip function properties
            if (typeof map[key] !== 'function') {
                keyValuePairs.push({ key: key, value: map[key] })
            }
        })
        return keyValuePairs
    }
});