import {fetchSubscriptions} from "./api/subscriptions";
import {bootstrap} from "./shared/common";
import {renderSubscriptionList} from "./domain/subscription-list";

Promise.all([
    bootstrap(),
    init(),
]);

async function init () {
    const [subscriptions, error] = await fetchSubscriptions();

    if (error) {
        console.error(error.message);
        return;
    }

    renderSubscriptionList('.js-subscriptions-list', subscriptions)
}
