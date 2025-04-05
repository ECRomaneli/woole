app.component('Dashboard', {
    template: /*html*/ `
        <session-details-card class="w-100" :session-details="sessionDetails"></session-details-card>
    `,
    props: { sessionDetails: Object }
})

app.component('SessionDetailsCard', {
    template: /*html*/ `
        <box label="Session Details" label-img="dashboard">
            <template #body>
                <table class="table table-striped table-hover" aria-label="Session Details">
                    <tbody>
                        <tr v-for="(value, key) in sessionDetails">
                            <template v-if='value'>
                                <td class="highlight">{{ getKey(key) }}</td>
                                <td v-html="getValue(value)"></td>
                            </template>
                        </tr>
                    </tbody>
                </table>
            </template>
        </box>
    `,
    props: { sessionDetails: Object },
    data() {
        return {
            keyMap: {
                clientId: 'Client ID',
                http: 'URL',
                https: 'Secure URL',
                proxying: 'Proxying',
                sniffer: 'Sniffer',
                tunnel: 'Tunnel URL',
                maxRecords: 'Max Stored Records'
            }
        }
    },
    methods: {
        getKey(key) {
            return this.keyMap[key] || key
        },

        getValue(value) {
            if ((value + "").indexOf("://") !== -1) {
                return '<a target="_blank" href="' + value + '">' + value + '</a>'
            }

            return value
        }
    }
})