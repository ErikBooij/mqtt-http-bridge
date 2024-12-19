import {bootstrap} from "./shared/common";
import {confirm, ConfirmStyle} from "./dom/confirm";
import {deleteGlobalParameter, fetchGlobalParameters} from "./api/global-parameters";

Promise.all([
    bootstrap(),
    init(),
]);

async function init () {

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
