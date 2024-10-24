import {bootstrap} from "./shared/common";
import {fetchSubscription} from "./api/subscriptions";
import {subscriptionForm} from "./forms/subscription-form";
import {fetchGlobalParameters} from "./api/global-parameters";
import {globalParameterForm} from "./forms/global-parameter-form";

Promise.all([
    bootstrap(),
    init(),
]);

declare global {
    interface Window {
        globalParameterKey: string;
    }
}

async function init(): Promise<void> {
    let oldValue = await (async() => {
        if (!window.globalParameterKey) {
            return undefined
        }

        const [params, error] = await fetchGlobalParameters();

        if (error !== null || params === null) {
            return undefined;
        }

        return params[window.globalParameterKey] ?? undefined;
    })();

    const form = globalParameterForm({
        key: window.globalParameterKey,
        selector: '.js-form',
        value: oldValue,
    });

    document.querySelector('.js-submit')?.addEventListener('click', (event) => {
        event.preventDefault();
        form.submit();
    });
}
