app.component('Dashboard', {
    template: /*html*/ `
        <div class="container-fluid">
            <div class="row">
                <session-details-card class="col-xl-6 col-lg-8 col-md-12" :session-details="sessionDetails"></session-details-card>
            </div>
        </div>
    `,
    props: { sessionDetails: Object }
})

app.component('SessionDetailsCard', {
    template: /*html*/ `
        <div class="p-0">
            <div class="card card-shadow" style="padding: 10px 20px">
                <div class="d-flex pt-3 pb-2"><span class="h2">Session Details</span></div>
                <table class="table table-striped table-hover" aria-label="Session Details">
                    <tbody>
                        <tr v-for="(value, key) in sessionDetails">
                            <th>{{ key }}</th>
                            <td v-html="getValue(value)"></td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
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