import {bootstrap} from "./shared/common";
import {fetchSubscription} from "./api/subscriptions";
import {subscriptionForm} from "./forms/subscription-form";
import {fetchGlobalParameters} from "./api/global-parameters";

Promise.all([
    bootstrap(),
    init(),
]);

declare global {
    interface Window {
        subscriptionID: string;
    }
}

async function init(): Promise<void> {
    const [
        [subscription, subscriptionError],
        [globalParameters, globalParametersError],
    ] = await Promise.all([
        (async () => {
            if (!window.subscriptionID) {
                return [undefined, null];
            }

            return fetchSubscription(window.subscriptionID);
        })(),
        fetchGlobalParameters(),
    ])

    if (subscriptionError !== null) {
        console.error(subscriptionError);
        return;
    }

    if (globalParametersError !== null) {
        console.error(globalParametersError);
        return;
    }

    const params = parameters(globalParameters);

    const form = subscriptionForm({
        selector: '.js-form',
        selectorParams: '.js-params',
        subscription,
        params,
    });

    document.querySelector('.js-submit')?.addEventListener('click', (event) => {
        event.preventDefault();
        form.submit();
    });
}

function parameters(globalParameters: Record<string, string>): Record<string, string[]> {
    return {
        'meta': ['topic', 'payload', 'client'],
        'global': Object.keys(globalParameters),
    };
}
