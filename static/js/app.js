// https://npm.reversehttp.com/#preact,preact/hooks,htm
// import {...} from 'https://npm.reversehttp.com/preact,preact/hooks,htm';
import { htm, h, useState, useEffect, useRef, render } from '/static/js/modules.js';

const baseHostname = 'pomegranate-winspector.haematite.dev';
const baseEndpoint = `https://${baseEndpoint}`;

const html = htm.bind(h);

const App = ({ initialBinKey }) => {
    const [binKey, setBinKey] = useState(initialBinKey);
    const [recordKey, setRecordKey] = useState(null);
    const [binData, setBinData] = useState({records: []});
    const ws = useRef(null);

    if (binKey === null) {
        return null;
    }

    const reloadBin = async () => {
        const r = await fetch(`/bin/${binKey}`, { method: "GET" });
        const data = await r.json();
        setBinData(data);
        return data;
    };

    const newBin = async () => {
        const r = await fetch("/bin", { method: "POST" });
        const data = await r.json();
        setRecordKey(null);
        setBinKey(data.bin_key);
        setBinData(data);
    };

    useEffect(async () => {
        const data = await reloadBin();
        const { records } = data;
        if (records.length > 0) {
            setRecordKey(records[0].record_key);
        }
    }, [binKey]);

    useEffect(async () => {
        if (ws.current !== null) {
            ws.current.close();
        }
        const socket = new WebSocket(`ws://${baseHostname}/bin/${binKey}/watch`);
        ws.current = socket;
        socket.onopen = () =>  { console.log("[open] connection established"); };
        socket.onclose = (event) => {
            if (event.wasClean) {
                console.log(`[close] connection closed cleanly, code=${event.code} reason=${event.reason}`);
            } else {
                // e.g. server process killed or network down
                // event.code is usually 1006 in this case
                console.log('[close] connection died');
            }
        };
        socket.onerror = (error) => { console.log(`[error] ${error.message}`); };
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
            records = [ record, ...records ];
            setBinData({ ...binData, records });
        };
    }, [binData]);

    const { records } = binData;
    const record = recordKey === null ? null : records.find(r => r.record_key == recordKey);

    return html`
        <${NavBar}
            binKey=${binKey}
            setBinKey=${setBinKey}
            reloadBin=${reloadBin}
            newBin=${newBin}
        />
        <${Content}
            records=${records}
            record=${record}
            setRecordKey=${setRecordKey}
        />
    `;
};


const NavBar = ({ binKey, setBinKey, reloadBin, newBin }) => {
    const onInput = (e) => { setBinKey(e.target.value); };
    return html`
        <header class="bg-dark text-white p-3">
            <div class="justify-content-center justify-content-lg-start align-items-center d-flex flex-wrap">
                <ul class="col-12 col-lg-auto nav justify-content-center me-lg-auto mb-2 mb-md-0">
                    <li class="d-md-none d-lg-block">
                        <a href="https://github.com/hellupline/winspector" class="nav-link px-2 text-white"> Winspector </a>
                    </li>
                </ul>
                <form class="col-12 col-lg-auto mb-3 mb-lg-0 me-lg-3">
                    <input type="text" class="form-control form-control-dark" placeholder="Bin Key..." value=${binKey} onInput=${onInput} />
                </form>
                <div class="text-end">
                    <${CopyRecorderUrlButton} binKey=${binKey} />
                    <${OpenInNewTabButton} binKey=${binKey} />
                    <${ReloadBinButton} reloadBin=${reloadBin} />
                    <${NewBinButton} newBin=${newBin} />
                </div>
            </div>
        </header>
    `;
};


const CopyRecorderUrlButton = ({ binKey }) => {
    if (binKey === null) {
        return null;
    }
    const onClick = (e) => {
        e.preventDefault();
        navigator.clipboard.writeText(`${baseEndpoint}/record/${binKey}`);
    };
    return html`<button class="btn btn-outline-light me-2" onClick=${onClick}> Copy url </button>`;
};


const OpenInNewTabButton = ({ binKey }) => {
    if (binKey === null) {
        return null;
    }
    return html`
        <a href="${baseEndpoint}/record/${binKey}" target="_blank" class="btn btn-outline-light me-2">
            Open in a new tab
        </a>
    `;
};


const ReloadBinButton = ({ reloadBin }) => {
    const onClick = (e) => { e.preventDefault(); reloadBin(); };
    return html`<button class="btn btn-outline-light me-2" onClick=${onClick}> Reload </button>`;
};


const NewBinButton = ({ newBin }) => {
    const onClick = (e) => { e.preventDefault(); newBin(); };
    return html`<button class="btn btn-warning" onClick=${onClick}> New Bin </button>`;
};


const Content = ({ records, record, setRecordKey }) => {
    return html`
        <div class="container-fluid">
            <div class="row justify-content-start">
                <div class="col-lg-2 col-sm-4">
                    <${RecordList} records=${records} setRecordKey=${setRecordKey} />
                </div>
                <div class="col-lg-10 col-sm-8">
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
    const { record_key: recordKey, created_at: createdAt, method } = record;
    const d = new Date(createdAt);
    const date = d.toLocaleDateString('ja-JP', { year: 'numeric', month: '2-digit', day: '2-digit' });
    const time = d.toLocaleTimeString('ja-JP', { hour12: false });
    const onClick = (e) => { e.preventDefault(); setRecordKey(recordKey); }
    return html`
        <a href="#" class="list-group-item list-group-item-action lh-tight py-3" onClick=${onClick}>
            <div class="justify-content-between align-items-center d-flex w-100">
                <strong class="text-truncate mb-1">
                    <span class="badge ${httpMethodColor(method)}">
                        ${method}
                    </span> - #${recordKey}
                </strong>
            </div>
            <div class="col-10 small mb-1"> ${date} ${time} </div>
        </a>
    `;
};


const RecordDetail = ({ record }) => {
    if (record === null) {
        return null;
    }
    const { headers, query, post_form_data: postFormData, body } = record;
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


const RequestTable = ({ record }) => {
    const {
        record_key: recordKey,
        created_at: createdAt,
        method,
        url,
        remote_addr: remoteAddr,
        content_lenght: contentLenght,
    } = record;
    const d = new Date(createdAt);
    const date = d.toLocaleDateString('ja-JP', { year: 'numeric', month: '2-digit', day: '2-digit' });
    const time = d.toLocaleTimeString('ja-JP', { hour12: false });
    return html`
        <h5>request details</h5>
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
                        <td class="text-nowrap"> size </td>
                        <td class="text-truncate"> ${contentLenght} bytes </td>
                    </tr>
                    <tr>
                        <td class="text-nowrap"> remote address </td>
                        <td class="text-truncate"> ${remoteAddr} </td>
                    </tr>
                    <tr>
                        <td class="text-nowrap"> id </td>
                        <td class="text-truncate"> ${date} ${recordKey} </td>
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
        <h5>${title}</h5>
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
        <h5>body</h5>
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


render(html`<${App} initialBinKey="d45a2464-4bce-4628-95be-8b8dfebe90be" />`, document.body);
