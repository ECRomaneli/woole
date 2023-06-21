app.component('Dashboard', {
    template: /*html*/ `
        <session-details-card class="w-100" :session-details="sessionDetails"></session-details-card>
    `,
    props: { sessionDetails: Object }
})

app.component('SessionDetailsCard', {
    template: /*html*/ `
        <box label="Session Datails" label-img="dashboard">
            <template #body>
                <table class="table table-striped table-hover" aria-label="Session Details">
                    <tbody>
                        <tr v-for="(value, key) in sessionDetails">
                            <td class="highlight">{{ key }}</td>
                            <td v-html="getValue(value)"></td>
                        </tr>
                    </tbody>
                </table>
            </template>
        </box>
    `,
    props: { sessionDetails: Object },
    methods: {
        getValue(value) {
            if ((value + "").indexOf("://") !== -1) {
                return '<a target="_blank" href="' + value + '">' + value + '</a>'
            }

            return value
        }
    }
})