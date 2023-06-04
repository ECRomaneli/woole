app.component('HeaderEditor', {
    template: /*html*/ `
        <table class="table header-table table-striped table-hover" aria-label="header">
            <thead>
                <tr>
                    <th v-if="!readOnly" scope="col" role="column:remove"></th>
                    <th scope="col" role="column:name">Name</th>
                    <th scope="col" role="column:value">Value</th>
                </tr>
            </thead>
            <tbody v-if="readOnly">
                <tr v-for="(header, index) in parsedHeader" :key="index">
                    <td class="highlight" role="name">{{ header.name }}</td>
                    <td role="value">{{ header.value }}</td>
                </tr>
            </tbody>
            <tbody v-else>
                <tr v-for="(header, index) in parsedHeader" :key="index">
                    <td><div class="c-pointer" @click="removeHeader(index)"><img class="svg-icon square-24" :src="$image.src('trash')" alt="remove header"></div></td>
                    <td><textarea placeholder="Name" class="auto-resize" spellcheck="false" @focus="autoResize" @input="autoResize" @blur="onBlur($event, header)" v-model="header.name"></textarea></td>
                    <td><textarea placeholder="Value" class="auto-resize" spellcheck="false" @focus="autoResize" @input="autoResize" @blur="onBlur($event, header)" v-model="header.value"></textarea></td>
                </tr>
                <tr>
                    <td colspan='3'><div class="c-pointer" @click="addHeader()"><img class="svg-icon square-24" :src="$image.src('plus')" alt="add header"></div></td>
                </tr>
            </tbody>
        </table>
    `,
    emits: [ 'update', 'remove' ],
    inject: ['$image'],
    props: { header: Object, readOnly: { type: Boolean, default: true } },
    data() { return { parsedHeader: this.parseHeader(this.header) } },
    watch: { header(header) { this.parsedHeader = this.parseHeader(header) } },

    methods: {
        addHeader() {
            this.parsedHeader.push({ name: '', value: '' })
        },

        removeHeader(index) {
            this.$emit('remove', this.parsedHeader.splice(index, 1)[0])
        },

        autoResize(event) {
            let el = event.currentTarget
            el.style.height = 'auto'
            el.style.height = event.type !== 'blur' ? el.scrollHeight + 'px' : ''
        },

        onBlur(event, header) {
            this.autoResize(event)
            this.$emit('update', header)
        },

        parseHeader(header) {
            const parsedHeader = []

            if (header === void 0) {
                console.warn("header should not be undefined")
                return parsedHeader
            }

            Object.keys(header).forEach(headerName => 
                parsedHeader.push({
                    name: headerName,
                    value: header[headerName].Val.join(';')
                })
            )
            return parsedHeader
        },

        getHeader() {
            let header = {}
            this.parsedHeader.forEach(h => header[h.name] = { Val: [h.value] })
            return header
        }
    }
})