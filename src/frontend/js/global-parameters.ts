import {ColumnWidthOption, Table, table} from "./dom/table";
import {bootstrap} from "./shared/common";
import {elementFromTemplate} from "./dom/template";
import {confirm, ConfirmStyle} from "./dom/confirm";
import {copyIcon} from "./shared/icons";
import {deleteGlobalParameter, fetchGlobalParameters} from "./api/global-parameters";

Promise.all([
    bootstrap(),
    init(),
]);

type GlobalParameterTableEntry = { key: string, value: string };

async function init () {
    const globalParametersTable = table<GlobalParameterTableEntry>('table.js-global-parameters-table', {
        order: ['key', 'value'],
        actions: ({ key }) => [
            { type: 'normal', label: 'Edit', href: '/global-parameters/' + key },
            { type: 'destructive', label: 'Delete', action: async (_: GlobalParameterTableEntry, t: Table<GlobalParameterTableEntry>) => { if (await handleDeleteGlobalParameter(key)) { t.removeRow(({ key: paramKey }) => paramKey === key ) } } },
        ],
        columns: {
            key: {
                label: 'Key',
                renderer: ({ key }) => {
                    const completeKey = `global.${key}`;

                    const copyButton = copyIcon();

                    copyButton.addEventListener('click', () => {
                        navigator.clipboard.writeText(completeKey)
                    })

                    const keyElement = elementFromTemplate(`<span title="${completeKey}" class="font-mono flex items-center gap-x-1">${completeKey}</span>`)

                    keyElement.append(copyButton)

                    return keyElement
                },
                width: ColumnWidthOption.MINIMAL,
            },
            value: {
                label: 'Value',
            },
        }
    });

    globalParametersTable.showLoading();

    const [globalParameters, error] = await fetchGlobalParameters();

    if (error) {
        globalParametersTable.showError(error.message);
        return;
    }


    globalParametersTable.showData(Object.entries(globalParameters).map(([key, value]) => ({ key, value })));
}

async function handleDeleteGlobalParameter(key: string): Promise<boolean> {
    return confirm({
        title: 'Delete Global Parameter',
        message: 'Are you sure you want to delete this global parameter?',
        style: ConfirmStyle.DESTRUCTIVE,
        confirmLabel: 'Delete',
        onConfirm: async () => {
            await deleteGlobalParameter(key)
        },
    })
}
