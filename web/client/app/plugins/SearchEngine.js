app.use({
    install: (app, _) => {
        const   TOKEN_SEPARATOR = ':', 
                REGEX_CHAR = '*', 
                RANGE_CHAR = '~',
                LEFT_PARENTHESIS_CHAR = '(',
                RIGHT_PARENTHESIS_CHAR = ')',
                EMPTY_QUOTES_STR = '""',
                KEY_SEPARATOR = '.',
                NEGATED_PREFIX = 'not ',
                RANGE_REGEXP = /^[-\D]*(-?\d+(\.\d+)?)?[-\D]*-[-\D]*(-?\d+(\.\d+)?)?[-\D]*$/,
                TOKENIZER = new RegExp(` *(${NEGATED_PREFIX}+)?([\\w${KEY_SEPARATOR}-]+)? *([${REGEX_CHAR}${RANGE_CHAR}]?${TOKEN_SEPARATOR})? *("((?:\\\\.|[^"\\\\])+)"|(?:\\\\.|[^ ()\\\\])+)? *(and|or|[()]|$)`, 'g'),
                TOKEN = { NEGATED: 1, KEY: 2, TYPE: 3, VALUE: 4, QUOTED_VALUE:5, QUALIFIER: 6 }
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

            return [...evaluateConditions(objList, extractConditionsFromQuery(queryStr.toLowerCase()), exclude)]
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

        function getQuery(negated, type, key, value) {
            const query = new Query()
            query.negated = negated
            query.key = key
            query.type = type
            query.value = value

            if (query.type === void 0) { return query }
            
            if (query.value === void 0 || query.value === null || query.value.trim() === EMPTY_STR) {
                delete query.type
                delete query.value
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

        function extractConditionsFromQuery(query, regex = new RegExp(TOKENIZER), group = new GroupQuery()) {
            while ((m = regex.exec(query)) !== null && m[0] !== EMPTY_STR) {
                const Q = m[TOKEN.QUALIFIER]

                if (Q === LEFT_PARENTHESIS_CHAR) {
                    group.addCondition(extractConditionsFromQuery(query, regex))
                    continue
                }

                let value = m[TOKEN.QUOTED_VALUE] || (m[TOKEN.VALUE] !== EMPTY_QUOTES_STR ? m[TOKEN.VALUE] : void 0)

                if (m[TOKEN.KEY] || value) {
                    const type = m[TOKEN.TYPE] && m[TOKEN.TYPE] !== TOKEN_SEPARATOR ? m[TOKEN.TYPE].charAt(0) : void 0
                    group.addCondition(getQuery(!!m[TOKEN.NEGATED], type, m[TOKEN.KEY], value))
                }

                if (Q === RIGHT_PARENTHESIS_CHAR) { break }

                if (Q) { group.addOperator(Operator.from(Q)) }
            }

            return group
        }

        function evaluateConditions(objList, group, exclude) {
            if (group.conditions.length === 0) { return objList }
            
            let currentResults = evaluateCondition(objList, group.conditions[0], exclude)
            
            for (let i = 1; i < group.conditions.length; i++) {
                const condition = group.conditions[i]
                const previousOperator = group.conditions[i - 1].operator
                
                if (previousOperator && previousOperator === Operator.OR) {
                    const nextResults = evaluateCondition(objList, condition, exclude)
                    currentResults.add(...nextResults)
                } else {
                    currentResults = evaluateCondition(currentResults, condition, exclude)
                }
            }
            
            return currentResults
        }
        
        function evaluateCondition(objList, condition, exclude) {
            if (condition.conditions) { return evaluateConditions(objList, condition, exclude) }
            const resultSet = new Set()
            objList.forEach(obj => {
                condition.negated !== findQuery(obj, condition, EMPTY_STR, exclude) && resultSet.add(obj)
            })
            return resultSet
        }
        

        class Query {
            constructor() {
                this.key = void 0
                this.value = void 0
                this.type = void 0
                this.operator = void 0
                this.negated = false
            }
        }

        class GroupQuery {
            constructor() {
                this.conditions = []
            }

            addCondition(condition) {
                this.conditions.push(condition)
            }

            addOperator(operator) {
                if (this.conditions.length === 0) { throw new Error('No conditions to add operator') }
                this.conditions[this.conditions.length-1].operator = operator
            }
        }

        class Operator {
            static AND = new Operator('and')
            static OR = new Operator('or')

            constructor(op) { this.op = op }

            static from(operator) {
                switch (operator) {
                    case this.OR.op: return this.OR
                    default: return this.AND
                }
            }
        }

        app.provide('$search', search)
        document.searchEngine = search
    }
})