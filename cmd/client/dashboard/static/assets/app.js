var app = new Vue({
    el: '#app',
    data: {
        items: [],
        selectedItem: {},
        proxyPort: '',
        targetURL: '',
        maxRecords: 16
    },
    created() {
        this.setupStream();
    },
    methods: {
        setupStream() {
            let es = new EventSource('/conn');
            es.addEventListener('config', event => {
                const cfg = JSON.parse(event.data);
                this.proxyPort = cfg.ProxyPort;
                this.ProxyURL = cfg.ProxyURL;
                this.maxRecords = cfg.MaxRecords;
            });
            es.addEventListener('records', event => {
                let items = JSON.parse(event.data) || [];
                this.items = items.reverse();
            });
            es.addEventListener('record', event => {
                let item = JSON.parse(event.data);
                if (item !== null) {
                    this.items.unshift(item);
                    while (this.items.length > this.maxRecords) {
                        this.items.pop();
                    }
                } else {
                    this.items = []
                }
            });
            es.onerror = () => {
                this.items = [];
                this.selectedItem = {};
            };
        },
        async show(item) {
            this.selectedItem = { ...this.selectedItem, id: item.id, status: item.status };
            let resp = await fetch('/info/' + item.id);
            let data = await resp.json();
            this.selectedItem = { ...this.selectedItem, ...data };
        },
        statusColor(item) {
            if (item.status < 300) return 'ok';
            if (item.status < 400) return 'warn';
            return 'error';
        },
        async clearDashboard() {
            this.selectedItem = {};
            await fetch('/clear');
        },
        canPrettifyBody(name) {
            if (!this.selectedItem[name]) return false;
            return this.selectedItem[name].indexOf('Content-Type: application/json') != -1;
        },
        prettifyBody(key) {
            let regex = /\n([\{\[](.*\s*)*[\}\]])/;
            let data = this.selectedItem[key];
            let match = regex.exec(data);
            let body = match[1];
            let prettyBody = JSON.stringify(JSON.parse(body), null, '    ');
            this.selectedItem[key] = data.replace(body, prettyBody);
        },
        copyCurl(event) {
            this.changeText(event);
            let e = document.createElement('textarea');
            e.value = this.selectedItem.curl;
            document.body.appendChild(e);
            e.select();
            document.execCommand('copy');
            document.body.removeChild(e);
        },
        async retry() {
            await fetch('/retry/' + this.selectedItem.id,
                { headers: { 'Cache-Control': 'no-cache' } });
            this.show(this.items[0]);
        },
        changeText(event) {
            let elem = event.target;
            let btnText = elem.getAttribute("data-text");
            elem.innerText = "copied!";
            setTimeout(() => elem.innerText = btnText, 400)
        }
    },
});