// https://npm.reversehttp.com/#preact,preact/hooks,preact-router,htm,history
// import {...} from 'https://npm.reversehttp.com/preact,preact/hooks,preact-router,htm';

import {
    render,
    h,
    useState,
    useEffect,
    useRef,
    Router,
    route,
    htm,
} from '/static/js/modules.js';


const baseEndpoint = `${window.location.protocol}//${window.location.host}`;
const baseSocketEndpoint = `${window.location.protocol === 'https:' ? 'wss:': 'ws:' }//${window.location.host}`;

const html = htm.bind(h);

const createBin = async () => {
    const r = await fetch("/bin", { method: "POST" });
    return await r.json();
};

const loadBin = async (binKey) => {
    const r = await fetch(`/bin/${binKey}`, { method: "GET" });
    return await r.json();
};


const websocketConnect = (binKey) => {
    const socket = new WebSocket(`${baseSocketEndpoint}/bin/${binKey}/watch`);
    socket.onopen = () =>  { console.log("[open] connection established"); };
    socket.onerror = (error) => { console.log(`[error] ${error.message}`); };
    socket.onclose = (event) => {
        if (event.wasClean) {
            console.log(`[close] connection closed cleanly, code=${event.code} reason=${event.reason}`);
        } else {
            // e.g. server process killed or network down
            // event.code is usually 1006 in this case
            console.log('[close] connection died');
        }
    };
    return socket;
};


const Root = () => {
    const [binData, setBinData] = useState(null);

    const newBin = async () => {
        const data = await createBin();
        setBinData(null);
        route(`/app/${data.bin_key}`);
    };

    const reloadBin = async () => {
        if (binData === null) {
            return;
        }
        const { bin_key: binKey } = binData;
        const data = await loadBin(binKey);
        setBinData(data);
    };

    return html`
        <${NavBar} binData=${binData} newBin=${newBin} reloadBin=${reloadBin} />
        <${Router}>
            <${App} path="/app/:binKey/:recordKey?" binData=${binData} setBinData=${setBinData} />
            <div default></div>
        </Router>
    `;
};


const NavBar = ({ binData, newBin, reloadBin }) => {
    return html`
        <header class="bg-dark text-white p-3">
            <div class="justify-content-center justify-content-lg-start align-items-center d-flex flex-wrap">
                <ul class="col-12 col-lg-auto nav justify-content-center me-lg-auto mb-2 mb-md-0">
                    <li class="d-md-none d-sm-none d-none d-lg-block">
                        <a href="/" class="nav-link px-2 text-white">
                            wInspector
                        </a>
                    </li>
                    <li class="d-md-none d-sm-none d-none d-lg-block">
                        <a href="https://github.com/hellupline/winspector" class="nav-link px-2 text-white">
                            github
                        </a>
                    </li>
                </ul>
                <div class="text-end">
                    <${Router}>
                        <${BinButtons} path="/app/:binKey/:recordKey?" reloadBin=${reloadBin} />
                    </Router>
                    <${NewBinButton} newBin=${newBin} />
                </div>
            </div>
        </header>
    `;
};

const BinButtons = ({ binKey, reloadBin }) => {
    return html`
        <a href="/record/${binKey}" target="_blank" class="btn btn-outline-light me-2">
            Open in a new tab
        </a>
        <${CopyRecorderUrlButton} binKey=${binKey} />
        <${ReloadBinButton} reloadBin=${reloadBin} />
    `;
};


const CopyRecorderUrlButton = ({ binKey }) => {
    const onClick = (e) => {
        e.preventDefault();
        navigator.clipboard.writeText(`${baseEndpoint}/record/${binKey}`);
    };
    return html`<button class="btn btn-outline-light me-2" onClick=${onClick}> Copy url </button>`;
};


const ReloadBinButton = ({ reloadBin }) => {
    const onClick = (e) => { e.preventDefault(); reloadBin(); };
    return html`<button class="btn btn-outline-light me-2" onClick=${onClick}> Reload </button>`;
};


const NewBinButton = ({ newBin }) => {
    const onClick = (e) => { e.preventDefault(); newBin(); };
    return html`<button class="btn btn-warning" onClick=${onClick}> New </button>`;
};


const App = ({ binKey, recordKey, binData, setBinData }) => {
    const loadBinData = async (binKey) => {
        const r = await fetch(`/bin/${binKey}`, { method: "GET" });
        const data = await r.json();
        setBinData(data);
    };

    const ws = useRef(null);

    useEffect(async () => {
        if (ws.current !== null) { ws.current.close(); }
        loadBinData(binKey);
        const socket = websocketConnect(binKey);
        ws.current = socket;
        return () => { socket.close(); }
    }, [binKey]);

    useEffect(async () => {
        if (ws.current === null) {
            return;
        }
        ws.current.onmessage = (event) => {
            console.log(`[message] data received from server: ${event.data}`);
            const record = JSON.parse(event.data);
            let { records } = binData;
            const length = records.length;
            records = [ record, ...records ];
            setBinData({ ...binData, records });
            console.log(`length = ${length}`);
            if (length === 0) {
                route(`/app/${binKey}/${record.record_key}`);
            }
        };
    }, [binData]);

    if (binData === null) {
        return null;
    }
    const { records } = binData;
    const record = recordKey === "" ? null : records.find(r => r.record_key == recordKey);
    return html`
        <div class="container-fluid px-0">
            <div class="row justify-content-start">
                <div class="col-lg-2 col-sm-4">
                    <${RecordList} records=${records} />
                </div>
                <div class="col-lg-10 col-sm-8 p-3">
                    <${RecordDetail} record=${record} />
                </div>
            </div>
        </div>
    `;
};


const RecordList = ({ records, setRecordKey }) => {
    return html`
        <div class="flex-column flex-shrink-0 align-items-stretch bg-white d-flex">
            <div class="list-group list-group-flush border-bottom scrollarea">
                ${records.map((record) => {
                    return html`<${RecordListItem}
                        record=${record}
                        setRecordKey=${setRecordKey}
                    />`;
                })}
            </div>
        </div>
    `;
};


const RecordListItem = ({ record, setRecordKey }) => {
    const {
        bin_key: binKey,
        record_key: recordKey,
        created_at: createdAt,
        method,
    } = record;
    const d = new Date(createdAt);
    const date = d.toLocaleDateString('ja-JP', { year: 'numeric', month: '2-digit', day: '2-digit' });
    const time = d.toLocaleTimeString('ja-JP', { hour12: false });
    let className = "list-group-item list-group-item-action lh-tight py-3";
    if (window.location.pathname === `/app/${binKey}/${recordKey}`) {
        className = `${className} active`;
    }
    return html`
        <a activeClassName="active" href="/app/${binKey}/${recordKey}" class="${className}">
            <div class="justify-content-between align-items-center d-flex w-100">
                <strong class="text-truncate mb-1">
                    <span class="badge ${httpMethodColor(method)}">
                        ${method}
                    </span> #${recordKey}
                </strong>
            </div>
            <div class="col-10 small mb-1"> ${date} ${time} </div>
        </a>
    `;
};


const RecordDetail = ({ record }) => {
    if (record === null) {
        return html`<${RecordWaiter} />`;
    }
    const {
        headers,
        query,
        post_form_data: postFormData,
        body
    } = record;
    return html`
        <div class="container-fluid pt-3">
            <div class="row justify-content-center">
                <div class="col">
                    <${RequestTable} record=${record} />
                </div>
                <div class="col">
                    <${KeyValueTable} items=${headers} title="headers" />
                </div>
                <div class="col">
                    <${KeyValueTable} items=${query} title="query" />
                </div>
                <div class="col">
                    <${KeyValueTable} items=${postFormData} title="post form data" />
                </div>
            </div>
            <div class="row justify-content-center">
                <div class="col">
                    <${RequestBody} body=${body} />
                </div>
            </div>
        </div>
    `;
};


const RecordWaiter = () => {
    return html`
        <div class="spinner-grow" role="status">
            <span class="visually-hidden">Loading...</span>
        </div>
    `;
}


const RequestTable = ({ record }) => {
    const {
        record_key: recordKey,
        created_at: createdAt,
        method,
        url,
        host,
        remote_addr: remoteAddr,
        content_lenght: contentLenght,
    } = record;
    const d = new Date(createdAt);
    const date = d.toLocaleDateString('ja-JP', { year: 'numeric', month: '2-digit', day: '2-digit' });
    const time = d.toLocaleTimeString('ja-JP', { hour12: false });
    return html`
        <h5 class="detail-header">request details</h5>
        <hr />
        <div class="table-responsive">
            <table class="table table-striped table-hover table-sm">
                <tbody>
                    <tr>
                        <td class="text-nowrap">
                            <span class="badge ${httpMethodColor(method)}">
                                ${method}
                            </span>
                        </td>
                        <td class="text-truncate">
                            <a href="${url}">${url}</a>
                        </td>
                    </tr>
                    <tr>
                        <td class="text-nowrap"> date </td>
                        <td class="text-truncate"> ${date} ${time} </td>
                    </tr>
                    <tr>
                        <td class="text-nowrap"> host </td>
                        <td class="text-truncate"> ${host} </td>
                    </tr>
                    <tr>
                        <td class="text-nowrap"> remote address </td>
                        <td class="text-truncate"> ${remoteAddr} </td>
                    </tr>
                    <tr>
                        <td class="text-nowrap"> size </td>
                        <td class="text-truncate"> ${contentLenght} bytes </td>
                    </tr>
                    <tr>
                        <td class="text-nowrap"> id </td>
                        <td class="text-truncate"> ${recordKey} </td>
                    </tr>
                </tbody>
            </table>
        </div>
    `;
};


const KeyValueTable = ({ items, title }) => {
    if (items.length === 0) {
        return null;
    }
    return html`
        <h5 class="detail-header">${title}</h5>
        <hr />
        <div class="table-responsive">
            <table class="table table-striped table-hover table-sm">
                <tbody>
                    ${items.map(({ key, value }) => {
                        return html`
                            <tr>
                                <td class="text-nowrap"> ${key.toLowerCase()} </td>
                                <td class="text-truncate"> ${value.toLowerCase()} </td>
                            </tr>
                        `;
                    })}
                </tbody>
            </table>
        </div>
    `;
};


const RequestBody = ({ body }) => {
    if (body === "") {
        return null;
    }
    try {
        body = JSON.stringify(JSON.parse(body), null, 4);
    } catch(e) {
        console.log('request body is not json');
    }
    return html`
        <h5 class="detail-header">body</h5>
        <div class="border rounded-3 bg-light p-3">
            <pre class="m-0">
                ${body}
            </pre>
        </div>
    `;
};


const httpMethodColor = (method) => {
    switch (method.toLowerCase()) {
        case 'get':
            return 'bg-success text-white';
        case 'post':
            return 'bg-primary text-white';
        case 'put':
            return 'bg-warning text-dark';
        case 'patch':
            return 'bg-info text-dark';
        case 'delete':
            return 'bg-danger text-white';
        case 'options':
            return 'bg-secondary text-white';
        default:
            return 'bg-secondary text-white';
    }
};


render(html`<${Root} />`, document.body);
