import {createSubscription, Subscription, SubscriptionWithoutID, updateSubscription} from "../api/subscriptions";
import {form} from "./form";
import {navigateToSubscriptionsOverview} from "../shared/navigation";
import {elementFromTemplate} from "../dom/template";
import {copyIcon} from "../shared/icons";

interface SubscriptionForm {
    submit(): void;
}

type SubscriptionFormOptions = {
    params: Record<string, string[]>;
    selector: string;
    selectorParams?: string;
    subscription?: Subscription
}

export const subscriptionForm = (
    {params, selector, selectorParams, subscription}: SubscriptionFormOptions
): SubscriptionForm => {
    const element = document.querySelector(selector);

    if (!element) {
        throw new Error(`Element with selector ${selector} not found`);
    }

    let extractIndex = 0;
    let headerIndex = 0;

    let paramsElement: HTMLElement | null = null;

    if (subscription) {
        params.extract = Object.keys(subscription.extract || {});
    }

    if (selectorParams) {
        paramsElement = document.querySelector(selectorParams);

        if (!paramsElement) {
            throw new Error(`Element with selector ${selectorParams} not found`);
        }

        renderParams(params, paramsElement);
    }


    const subForm = form<SubscriptionWithoutID>({
        selector,
        transformer: (data: FormData): SubscriptionWithoutID => {
            return {
                name: data.get('name') as string,
                topic: data.get('topic') as string,
                extract: Object.fromEntries(
                    Array.from(data.getAll('extractVar[]'))
                        .map((key, index) => [key, data.getAll('extractValue[]')[index]])
                        .filter(([key, value]) => key !== '' && value !== '')
                ),
                filter: data.get('filter') as string,
                method: data.get('method') as Subscription['method'],
                url: data.get('url') as string,
                headers: Object.fromEntries(
                    Array.from(data.getAll('headerName[]'))
                        .map((key, index) => [key, data.getAll('headerValue[]')[index]])
                        .filter(([key, value]) => key !== '' && value !== '')
                ),
                body: data.get('body') as string,
            };
        },
        validator: (data: SubscriptionWithoutID) => {
            const errors = [];

            if (data.name.trim() === '') {
                errors.push({field: 'Name', message: 'Value is required'});
            }

            if (data.topic.trim() === '') {
                errors.push({field: 'Topic', message: 'Value is required'});
            }

            if (data.method.trim() === '') {
                errors.push({field: 'Method', message: 'Value is required'});
            }

            if (data.url.trim() === '') {
                errors.push({field: 'URL', message: 'Value is required'});
            }

            return errors;
        }
    });

    const addExtractField = (key: string, value: string) => {
        subForm.addFields(
            'extract',
            extractIndex++ === 0 ? 'Extract' : '',
            {name: 'extractVar[]', type: 'text', value: key, width: 'narrow'},
            {name: 'extractValue[]', type: 'text', value: value || '', mono: true},
        )
    }

    const addHeaderField = (key: string, value: string) => {
        subForm.addFields(
            'headers',
            headerIndex++ === 0 ? 'Headers' : '',
            {name: 'headerName[]', type: 'text', value: key || '', width: 'narrow'},
            {name: 'headerValue[]', type: 'text', value: value || '', mono: true},
        )
    }

    subForm.onSubmit(subscription
        ? async (data: SubscriptionWithoutID) => {
            console.log('Update subscription', data);

            const [_, error] = await updateSubscription({...data, id: subscription.id});

            console.log(error);

            if (error === null) {
                navigateToSubscriptionsOverview();
            }
        }
        : async (data: SubscriptionWithoutID) => {
            console.log('Create subscription', data);

            const [_, error] = await createSubscription(data);

            console.log(error);

            if (error === null) {
                navigateToSubscriptionsOverview();
            }
        });

    subForm.onChange((data) => {
        if (Object.entries(data.extract || {}).length === extractIndex) {
            addExtractField('', '');
        }

        if (Object.entries(data.headers || {}).length === headerIndex) {
            addHeaderField('', '');
        }

        const prevExtract = params.extract || [];

        params.extract = Object.keys(data.extract || {});

        if (paramsElement) {
            if (!prevExtract.every((key) => params.extract.includes(key)) || !params.extract.every((key) => prevExtract.includes(key))) {
                renderParams(params, paramsElement);
            }
        }
    });

    subForm
        .addSection({ name: 'identifier'})
        .addField({section: 'identifier', label: 'Name', name: 'name', type: 'text', value: subscription?.name || ''})
        .addSection({ name: 'preprocessing', label: 'Input', description: 'Configure which messages to act on, and how to extract data from them' })
        .addField({
            section: 'preprocessing',
            label: 'Topic',
            name: 'topic',
            type: 'text',
            value: subscription?.topic || ''
        })
        .addSection({name: 'extract', subsectionOf: 'preprocessing'});

    for (const [key, value] of Object.entries(subscription?.extract || {})) {
        addExtractField(key, value);
    }

    // Additional row for new extractions.
    addExtractField('', '');

    subForm.addField({
        section: 'preprocessing',
        label: 'Filter',
        name: 'filter',
        type: 'text',
        value: subscription?.filter || '',
        width: 'wide'
    })
        .addSection({ name: 'request', label: 'Output', description: 'The HTTP request to send when a message is received' })
        .addField({
            section: 'request',
            label: 'Method',
            name: 'method',
            type: 'select',
            options: ['GET', 'POST', 'PATCH', 'PUT', 'DELETE', 'HEAD', 'OPTIONS'],
            value: subscription?.method || 'POST',
            width: 'narrow'
        })
        .addField({section: 'request', label: 'URL', name: 'url', type: 'text', value: subscription?.url || ''})
        .addSection({name: 'headers', subsectionOf: 'request'});

    for (const [key, value] of Object.entries(subscription?.headers || {})) {
        addHeaderField(key, value);
    }

    // Additional row for new headers.
    addHeaderField('', '');

    subForm.addField({
        section: 'request',
        label: 'Body',
        name: 'body',
        type: 'textarea',
        value: subscription?.body || '',
        width: 'wide',
        mono: true
    });

    return {
        submit: () => {
            subForm.submit();
        }
    };
}

function renderParams (params: Record<string, string[]>, container: HTMLElement): void {
    container.classList.add('space-y-3', 'text-sm', 'text-gray-600');
    container.innerHTML = '';

    for (const [group, keys] of Object.entries(params).sort(([a], [b]) => a.localeCompare(b))) {
        const list = elementFromTemplate(`<ul class="overflow-hidden rounded-md bg-white px-4 py-4 shadow"></ul>`);

        for (const key of keys.sort()) {
            const item = elementFromTemplate(`<li class="whitespace-nowrap flex items-center gap-x-2 py-1"><span class="font-mono text-sm">{{ .${group}.<span class="font-bold">${key}</span> }}</span></li>`);

            const copyButton = copyIcon();
            copyButton.addEventListener('click', () => {
                navigator.clipboard.writeText(`{{ .${group}.${key} }}`);
            })

            item.prepend(copyButton);
            list.append(item);
        }

        container.append(list);
    }
}
