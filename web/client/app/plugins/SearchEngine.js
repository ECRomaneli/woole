app.use({
    install: (app, _) => {
        const   TOKEN_SEPARATOR = ':', 
                REGEX_CHAR = '*', 
                RANGE_CHAR = '~',
                KEY_SEPARATOR = '.',
                RANGE_REGEXP = new RegExp(`^[-\\D]*(-?\\d+(\\.\\d+)?)?[-\\D]*-[-\\D]*(-?\\d+(\\.\\d+)?)?[-\\D]*$`), 
                NOT_PREFIX = 'not ',
                UNKNOWN = -1, 
                EMPTY_ARR = [], 
                EMPTY_STR = '',
                STRING = 'string', 
                NUMBER = 'number', 
                BOOLEAN = 'boolean',
                BIGINT = 'bigint'
        
        function search(objList, queryStr, exclude) {
            if (objList === void 0 || objList === null) { return [] }
            if (queryStr === void 0 || queryStr === null || queryStr.trim() === EMPTY_STR) { return objList.slice() }
    
            const query = getQuery(queryStr)
            return objList.filter((obj) => query.not !== findQuery(obj, query, EMPTY_STR, exclude))
        }

        function findQuery(obj, query, nestedKeys, excludedKeys, keyFound) {
            return getObjectKeys(obj).some((key) => {
                const newNestedKeys = nestedKeys + KEY_SEPARATOR + key.toLowerCase()

                if (isExcluded(newNestedKeys, excludedKeys)) { return false }
                
                if (keyFound === void 0) {
                    if (newNestedKeys.indexOf(query.key) === UNKNOWN) {
                        return findQuery(obj[key], query, newNestedKeys, excludedKeys) 
                            || (query.value === void 0 && match(query.key, obj[key], query.type))
                    }

                    if (query.value === void 0) { return true }
                }

                return match(query.value, obj[key], query.type) || 
                       findQuery(obj[key], query, newNestedKeys, excludedKeys, true)
            })
        }

        function match(expectedValue, value, type) {
            if (value === null || value === void 0) { return false }
            
            const typeOf = typeof value

            // Range match
            if (type === RANGE_CHAR) {
                if (typeOf !== NUMBER && typeOf !== BIGINT) { return false }
                return matchRange(expectedValue, Number(value))
            }
            
            // Regex match
            if (type === REGEX_CHAR) { return expectedValue.test(value) }

            // Standard string match
            if (typeOf === STRING) {
                return `${value}`.toLowerCase().indexOf(expectedValue) !== UNKNOWN
            }

            // Convert numbers and booleans to string and check
            if (typeOf === NUMBER || typeOf === BIGINT || typeOf === BOOLEAN)  {
                return `${value}`.indexOf(expectedValue) !== UNKNOWN
            }

            return false
        }

        function matchRange(expectedRange, numValue) {
            if (expectedRange.min !== void 0 && expectedRange.max !== void 0) {
                return numValue >= expectedRange.min && numValue <= expectedRange.max
            } 
            if (expectedRange.min !== void 0) { return numValue >= expectedRange.min }
            return numValue <= expectedRange.max
        }

        function getQuery(rawQuery) {
            const query = { 
                not: rawQuery.indexOf(NOT_PREFIX) === 0
            }
            
            if (query.not) { rawQuery = rawQuery.substr(NOT_PREFIX.length) }

            const SEPARATOR_INDEX = rawQuery.indexOf(TOKEN_SEPARATOR)
            if (SEPARATOR_INDEX !== UNKNOWN) {
                const rawType = rawQuery[SEPARATOR_INDEX - 1]
                query.type = rawType === REGEX_CHAR ? REGEX_CHAR :
                             rawType === RANGE_CHAR ? RANGE_CHAR : void 0
            }

            const separator = query.type ? query.type + TOKEN_SEPARATOR : TOKEN_SEPARATOR
            const tokens = rawQuery.trim().split(separator)
            query.key = tokens.shift().trim().toLowerCase()
            query.value = tokens.join(separator).trim()

            if (query.value === void 0 || query.value === null || query.value.trim() === EMPTY_STR) {
                delete query.type
                delete query.value
                return query
            }

            if (query.type === void 0) {
                query.value = query.value.toLowerCase()
                return query
            }

            if (query.type === REGEX_CHAR) {
                try {
                    query.value = new RegExp(query.value, 'i')
                } catch (_e) {
                    delete query.type
                    delete query.value
                }
                return query
            }            
            
            const matches = query.value.match(RANGE_REGEXP)
            if (matches === null) {
                delete query.type
                delete query.value
                return query
            }

            query.value = { min: parseFloat(matches[1]) || void 0, max: parseFloat(matches[3]) || void 0 }

            if (query.value.min === void 0 && query.value.max === void 0) {
                delete query.type
                delete query.value
            }

            return query
        }
        
        function isExcluded(nestedKeys, excludedKeys) {
            return excludedKeys !== void 0 && excludedKeys.some((key) => nestedKeys.endsWith(key))
        }
        
        function getObjectKeys(obj) {
            return obj instanceof Object ? Object.keys(obj) : EMPTY_ARR
        }

        app.provide('$search', search)
        document.searchEngine = search
    }
})