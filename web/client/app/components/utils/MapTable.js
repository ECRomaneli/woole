app.component('MapTable', {
    template: /*html*/ `
        <table class="table map-table table-striped" :class="{ 'table-hover': readOnly }" aria-label="map table">
            <thead>
                <tr>
                    <th v-if="!readOnly" scope="col" role="column:remove"></th>
                    <th scope="col" role="column:key">Key</th>
                    <th scope="col" role="column:value">Value</th>
                </tr>
            </thead>
            <tbody v-if="readOnly">
                <tr v-for="(keyValuePair, index) in keyValuePairs" :key="index">
                    <td class="highlight" role="key">{{ keyValuePair.key }}</td>
                    <td role="value">{{ supress?.includes(keyValuePair.key) ? '...' : keyValuePair.value }}</td>
                </tr>
            </tbody>
            <tbody v-else>
                <tr v-for="(keyValuePair, index) in keyValuePairs" :key="index">
                    <td class="c-pointer" @click="removeKey(index)"><img class="svg-icon square-24" :src="$image.src('trash')" alt="remove key"></td>
                    <td><textarea name="key"   placeholder="Key"   class="auto-resize" spellcheck="false" @focus="autoResize" @input="autoResize" @blur="onBlur($event, keyValuePair)" v-model="keyValuePair.key"></textarea></td>
                    <td><textarea name="value" placeholder="Value" class="auto-resize" spellcheck="false" @focus="autoResize" @input="autoResize" @blur="onBlur($event, keyValuePair)" v-model="keyValuePair.value"></textarea></td>
                </tr>
                <tr class="c-pointer" @click="addKey()">
                    <td><img class="svg-icon square-24" :src="$image.src('plus')" alt="add key"></td>
                    <td></td>
                    <td></td>
                </tr>
            </tbody>
        </table>
    `,
    emits: [ 'update', 'remove' ],
    inject: ['$image', '$map'],
    props: { map: Object, readOnly: { type: Boolean, default: true }, supress: Array },
    data() { return { keyValuePairs: this.$map.toKeyValuePairs(this.map) } },
    watch: { map(newMap) { this.keyValuePairs = this.$map.toKeyValuePairs(newMap) } },
    methods: {
        addKey() { this.keyValuePairs.push({ key: '', value: '' }) },
        removeKey(index) { this.$emit('remove', this.keyValuePairs.splice(index, 1)[0]) },

        autoResize(event) {
            const el = event.currentTarget
            el.style.height = ''
            if (event.type !== 'blur') {
                el.style.height = Math.max(el.scrollHeight, el.offsetHeight) + 'px'
            }
        },

        onBlur(event, keyValuePair) {
            this.autoResize(event)
            this.$emit('update', keyValuePair)
        },

        toMap() {
            let map = {}
            this.keyValuePairs.forEach(p => map[p.key] = p.value)
            return map
        }
    }
})