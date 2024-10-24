import {deleteSubscription, fetchSubscriptions, Subscription} from "./api/subscriptions";
import {ColumnWidthOption, Table, table} from "./dom/table";
import {bootstrap} from "./shared/common";
import {elementFromTemplate} from "./dom/template";
import {confirm, ConfirmStyle} from "./dom/confirm";
import {copyIcon} from "./shared/icons";

Promise.all([
    bootstrap(),
    init(),
]);

type SubscriptionTableEntry = Subscription & { target?: string };

async function init () {
    const subscriptionsTable = table<SubscriptionTableEntry>('table.js-subscriptions-table', {
        order: ['id', 'name', 'topic', 'target'],
        actions: ({ id }) => [
            { type: 'normal', label: 'Edit', href: '/subscriptions/' + id },
            { type: 'destructive', label: 'Delete', action: async (_: SubscriptionTableEntry, t: Table<SubscriptionTableEntry>) => { if (await handleDeleteSubscription(id)) { t.removeRow(({ id: subId }) => subId === id ) } } },
        ],
        columns: {
            id: {
                label: 'ID',
                renderer: ({ id }) => {
                    const copyButton = copyIcon();

                    copyButton.addEventListener('click', () => {
                        navigator.clipboard.writeText(id)
                    })

                    const idElement = elementFromTemplate(`<span title="${id}" class="font-mono flex items-center gap-x-1">${id.substring(0, 8)}&hellip;</span>`)

                    idElement.append(copyButton)

                    return idElement
                },
                width: ColumnWidthOption.MINIMAL,
            },
            name: {
                label: 'Name',
            },
            topic: {
                label: 'Topic',
            },
            target: {
                label: 'Target',
                renderer: ({ method, url }) => {
                    return elementFromTemplate(`<span class="flex items-center gap-x-1"><span class="font-mono text-xs tracking-wider ${httpMethodColors[method]} font-bold px-1 py-0.5 rounded-sm">${method}</span> ${url}</span>`)
                },
            },
        }
    });

    subscriptionsTable.showLoading();

    const [subscriptions, error] = await fetchSubscriptions();

    if (error) {
        subscriptionsTable.showError(error.message);
        return;
    }


    subscriptionsTable.showData(subscriptions);
}

async function handleDeleteSubscription(id: string): Promise<boolean> {
    return confirm({
        title: 'Delete Subscription',
        message: 'Are you sure you want to delete this subscription?',
        style: ConfirmStyle.DESTRUCTIVE,
        confirmLabel: 'Delete',
        onConfirm: async () => {
            await deleteSubscription(id)
        },
    })
}

const httpMethodColors = {
    GET: 'bg-green-300',
    POST: 'bg-blue-300',
    PUT: 'bg-yellow-300',
    DELETE: 'bg-red-300',
    PATCH: 'bg-purple-300',
    OPTIONS: 'bg-gray-300',
    HEAD: 'bg-gray-300',
    TRACE: 'bg-gray-300',
    CONNECT: 'bg-gray-300',
};
