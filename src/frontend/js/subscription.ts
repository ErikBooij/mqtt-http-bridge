import {bootstrap} from "./shared/common";
import {elementFromTemplate} from "./dom/template";
import {copyIcon} from "./shared/icons";
import {fetchGlobalParameters} from "./api/global-parameters";
import {fetchSubscription, Subscription} from "./api/subscriptions";
import {validateExtract, validateFilter} from "./forms/subscription-validation";

Promise.all([
    bootstrap(),
    init(),
]);

declare global {
    interface Window {
        subscriptionForm: typeof subscriptionForm;
    }
}

async function init() {
    await subscriptionForm();
}

async function subscriptionForm(subscriptionId?: string, mode?: 'create' | 'edit') {
    const [
        [globalParameters, paramsError],
        [subscription, subscriptionError],
    ] = await Promise.all([
        fetchGlobalParameters(),
        subscriptionId ? fetchSubscription(subscriptionId) : [null, null],
    ]);

    if (paramsError) {
        console.error(paramsError.message);
        return;
    }

    if (subscriptionError) {
        console.error(subscriptionError.message);
        return;
    }

    const paramsElement = document.querySelector('.js-params');

    if (!paramsElement) {
        console.error('Params element not found');
        return;
    }

    const rerenderParams = (extractedParams: Record<string, string>) => {
        renderParams(parameters(globalParameters, extractedParams), paramsElement);
    }

    rerenderParams(subscription?.extract || {})

    const element = document.querySelector<HTMLFormElement>('.js-form');

    if (!element) {
        console.error('Form not found');
        return;
    }

    if (subscription) prepopulateForm(element, subscription);

    attachEventListeners(element, rerenderParams);
}

function addField(element: HTMLElement, type: 'extract' | 'header') {
    const tpl = ((): HTMLElement | null => {
        switch (type) {
            case 'extract':
                return element.querySelector('.js-extract-template');
            case 'header':
                return element.querySelector('.js-header-template');
        }

        return null
    })();

    if (!tpl) return;

    const duplicate = tpl.cloneNode(true) as HTMLElement;

    duplicate.querySelector('label')?.classList.add('invisible');

    tpl.parentNode?.append(duplicate);
}

function attachEventListener(element: HTMLElement | null, name: string, validation: Validation) {
    if (!element) {
        console.error(`Element with name ${name} not found`);
        return;
    }

    const errorElement = element.closest('*:has(.js-error)')?.querySelector('.js-error');

    if (!errorElement) {
        console.error(`Error element for ${name} not found`);
        return;
    }

    element.addEventListener('input', onInput(validation, errorElement as HTMLElement, name));
}

function attachEventListeners(form: HTMLFormElement, rerenderParams: (extractedParams: Record<string, string>) => void) {
    attachEventListener(form.querySelector('[name="name"]'), 'Name', { required: true });
    attachEventListener(form.querySelector('[name="topic"]'), 'Topic', { required: true });

    form.querySelectorAll('[name="extractVar[]"]').forEach((element) => {
        attachEventListener(element as HTMLElement, 'Extract Variable', {  });
    });
    form.querySelectorAll('[name="extractValue[]"]').forEach((element) => {
        attachEventListener(element as HTMLElement, 'Extract Value', { validExtract: true });
    });

    attachEventListener(form.querySelector('[name="filter"]'), 'Filter', { validFilter: true });
    attachEventListener(form.querySelector('[name="method"]'), 'Method', { required: true });
    attachEventListener(form.querySelector('[name="url"]'), 'URL', { required: true });

    form.querySelectorAll('[name="headerName[]"]').forEach((element) => {
        attachEventListener(element as HTMLElement, 'Header Name', {  });
    });
    form.querySelectorAll('[name="headerValue[]"]').forEach((element) => {
        attachEventListener(element as HTMLElement, 'Header Value', {  });
    });

    attachEventListener(form.querySelector('[name="body"]'), 'Body', {  });
}

function parameters(globalParameters: Record<string, string>, extractedParameters: Record<string, string>): Record<string, string[]> {
    return {
        'extract': Object.keys(extractedParameters),
        'meta': ['topic', 'payload', 'client'],
        'global': Object.keys(globalParameters),
    };
}

function prepopulateForm(element: HTMLFormElement, subscription: Subscription) {
    for (let i = 0; i < Object.entries(subscription.extract || {}).length; i++) {
        addField(element, 'extract');
    }

    for (let i = 0; i < Object.entries(subscription.headers || {}).length; i++) {
        addField(element, 'header');
    }

    setFormValue(element, 'name', subscription.name);
    setFormValue(element, 'topic', subscription.topic);
    setFormValues(element, 'extractVar', 'extractValue', subscription.extract);
    setFormValue(element, 'filter', subscription.filter);
    setFormValue(element, 'method', subscription.method);
    setFormValue(element, 'url', subscription.url);
    setFormValues(element, 'headerName', 'headerValue', subscription.headers);
    setFormValue(element, 'body', subscription.body);
}

function renderParams (params: Record<string, string[]>, container: Element): void {
    container.classList.add('space-y-3', 'text-sm', 'text-gray-600');
    container.innerHTML = '';

    for (const [group, keys] of Object.entries(params).sort(([a], [b]) => a.localeCompare(b))) {
        if (!keys.length) continue;

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

function setFormValue(form: HTMLFormElement, field: string, value: string | undefined, index?: number) {
    const element = index !== undefined ? form.querySelectorAll(`[name="${field}"]`)[index] : form.querySelector(`[name="${field}"]`);

    if (!element) {
        if (index !== undefined) {
            console.error(`Element with name ${field}[${index}] not found`);
        } else {
            console.error(`Element with name ${field} not found`);
        }
        return;
    }

    if (element.tagName.toLowerCase() === 'textarea') {
        element.textContent = value || '';
    } else {
        element.setAttribute('value', value || '');
    }
}

function setFormValues(form: HTMLFormElement, keyField: string, valueField: string, values: Record<string, string> | undefined) {
    if (!values) {
        return;
    }

    let i = 0;

    Object.entries(values).forEach(([key, value]) => {
        setFormValue(form, `${keyField}[]`, key, i);
        setFormValue(form, `${valueField}[]`, value, i);

        i++;
    })
}

type Validation = {
    required?: boolean;
    validExtract?: boolean;
    validFilter?: boolean;
}

function onInput(validation: Validation, errorElement: HTMLElement, field: string, index?: number) {
    const renderError = (message: string) => {
        errorElement.textContent = message;
    }

    const clearError = () => {
        errorElement.textContent = '';
    }

    let abortController: AbortController;

    return async function (event: Event) {
        if (abortController) {
            abortController.abort()
        }

        abortController = new AbortController();

        if (validation.required && (event.target as HTMLInputElement).value.trim() === '') {
            return renderError(`${field} is required`);
        }

        if (validation.validExtract && (event.target as HTMLInputElement).value.trim() !== '') {
            const error = await validateExtract(abortController, (event.target as HTMLInputElement).value)

            if (error) {
                return renderError(`${field}: ${error}`);
            }
        }

        if (validation.validFilter && (event.target as HTMLInputElement).value.trim() !== '') {
            const error = await validateFilter(abortController, (event.target as HTMLInputElement).value)

            if (error) {
                return renderError(`${field}: ${error}`);
            }
        }

        return clearError();
    }
}

window.subscriptionForm = subscriptionForm;
