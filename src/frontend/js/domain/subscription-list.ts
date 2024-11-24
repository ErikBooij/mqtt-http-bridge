import {deleteSubscription, Subscription} from "../api/subscriptions";
import {elementFromTemplate} from "../dom/template";
import {copyIcon} from "../shared/icons";
import {confirm, ConfirmStyle} from "../dom/confirm";

export const renderSubscriptionList = (selector: string, subscriptions: Subscription[]) => {
    const element = document.querySelector(selector);

    if (!element) {
        throw new Error(`Element with selector ${selector} not found`);
    }

    const list = elementFromTemplate(`<ul role="list" class="divide-y divide-gray-100"></ul>`);

    for (const subscription of subscriptions) {
        const item = elementFromTemplate(`
            <li class="flex items-center justify-between gap-x-6 py-5">
                <div class="min-w-0">
                    <div class="flex items-start gap-x-3">
                        <p class="text-sm/6 font-semibold text-gray-900">${subscription.name}</p>
                        <p class="mt-0.5 whitespace-nowrap rounded-md bg-green-50 px-1.5 py-0.5 text-xs font-medium text-green-700 ring-1 ring-inset ring-green-600/20">Active</p>
                    </div>
                    <div class="mt-1 flex items-center gap-x-2 text-xs/5 text-gray-500">
                        <p class="whitespace-nowrap js-subscription-id flex items-center gap-x-1 cursor-pointer">${subscription.id}</p>
                        <svg viewBox="0 0 2 2" class="h-0.5 w-0.5 fill-current">
                            <circle cx="1" cy="1" r="1" />
                        </svg>
                        <p class="truncate">${subscription.topic}</p>
                    </div>
                </div>
                <div class="flex flex-none items-center gap-x-4">
                    <a href="/subscriptions/${subscription.id}" class="hidden rounded-md bg-white px-2.5 py-1.5 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 sm:block">Edit</a>
                    <div class="relative flex-none">
                        <button type="button" class="-m-2.5 block p-2.5 text-gray-500 hover:text-gray-900 js-subscription-menu-btn">
                            <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true" data-slot="icon">
                                <path d="M10 3a1.5 1.5 0 1 1 0 3 1.5 1.5 0 0 1 0-3ZM10 8.5a1.5 1.5 0 1 1 0 3 1.5 1.5 0 0 1 0-3ZM11.5 15.5a1.5 1.5 0 1 0-3 0 1.5 1.5 0 0 0 3 0Z" />
                            </svg>
                        </button>
                        <div class="absolute right-0 z-10 mt-2 w-32 origin-top-right rounded-md bg-white py-2 shadow-lg ring-1 ring-gray-900/5 focus:outline-none hidden js-subscription-menu">
                            <a href="/subscriptions/${subscription.id}" class="block px-3 py-1 text-sm/6 text-gray-900 sm:hidden" role="menuitem" tabindex="-1" id="options-menu-0-item-0">Edit</a>
                            <a href="/new-subscription?base=${subscription.id}" class="block px-3 py-1 text-sm/6 text-gray-900" role="menuitem" tabindex="-1" id="options-menu-0-item-0">Duplicate</a>
                            <a href="#" class="block px-3 py-1 text-sm/6 text-red-900 js-delete-subscription" role="menuitem" tabindex="-1" id="options-menu-0-item-2">Delete</a>
                        </div>
                    </div>
                </div>
            </li>
        `);

        let menuIsOpen = false;
        const subscriptionMenuButton = item.querySelector('.js-subscription-menu-btn');
        const subscriptionMenu = item.querySelector('.js-subscription-menu');

        subscriptionMenuButton?.addEventListener('click', () => {
            menuIsOpen = !menuIsOpen;

            console.log({ menuIsOpen })

            if (menuIsOpen) {
                subscriptionMenu?.classList.remove('hidden');

                return
            }

            subscriptionMenu?.classList.add('hidden');
        });

        const idCopyButton = copyIcon();
        idCopyButton.title = 'Copy subscription ID';

        const subscriptionIdContainer = item.querySelector('.js-subscription-id')

        subscriptionIdContainer?.append(idCopyButton);
        subscriptionIdContainer?.addEventListener('click', () => {
            navigator.clipboard.writeText(subscription.id)
        })

        item.querySelector('.js-delete-subscription')?.addEventListener('click', async () => {
            const confirmed = await handleDeleteSubscription(subscription.id);

            if (confirmed) {
                item.remove();
            }
        })

        list.append(item)
    }

    element.innerHTML = '';
    element.append(list)
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
