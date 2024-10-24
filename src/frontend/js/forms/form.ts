import {elementFromTemplate} from "../dom/template";
import {DistributiveOmit} from "../shared/types";

export interface Form<T> {
    addField(options: AddFieldOptions): Form<T>;

    addFields(section: string, label: string, ...options: DistributiveOmit<AddFieldOptions, 'section' | 'label'>[]): Form<T>;

    addSection(options: AddSectionOptions): Form<T>;

    onChange(callback: (data: T) => void): Form<T>;

    onSubmit(callback: (data: T) => Promise<void> | void): Form<T>;

    submit(): void;
}

type FormOptions<T> = {
    selector: string;
    transformer: (data: FormData) => T;
    validator?: (data: T) => ValidationError[];
}

type ValidationError = { field: string; message: string; }

type AddFieldOptions = FieldOptions & { section: string };
type AddSectionOptions = { name: string } & ({ label?: never; description?: never; subsectionOf: string } | {
    label?: string | undefined;
    description?: string | undefined;
    subsectionOf?: never
});

export const form = <T>(options: FormOptions<T>): Form<T> => {
    const f = {} as Form<T>;
    const changeHandlers: Array<(data: T) => void> = [];
    const submitHandlers: Array<(data: T) => void> = [];

    const sections: { [key: string]: HTMLElement } = {};

    const element = document.querySelector(options.selector);

    if (!element) {
        throw new Error(`Element with selector ${options.selector} not found`);
    }

    const formElement = elementFromTemplate<HTMLFormElement>(`<form></form>`);

    formElement.classList.add('space-y-8', 'border-b', 'pb-12', 'sm:space-y-8', 'sm:divide-y', 'sm:divide-gray-900/10', 'sm:divide-y-2', 'sm:pb-16')

    element.append(formElement);

    const onChangeHandler = () => {
        const data = options.transformer(new FormData(formElement));

        changeHandlers.forEach(handler => handler(data));
    }

    f.addField = ({section, label, ...options}: AddFieldOptions): Form<T> => {
        return f.addFields(section, label, options)
    }

    f.addFields = (section: string, label: string, ...options: DistributiveOmit<AddFieldOptions, 'section' | 'label'>[]): Form<T> => {
        const sec = sections[section];

        if (!sec) {
            throw new Error(`Section with name ${section} not found`);
        }

        const row = createFieldRow(onChangeHandler, label, ...options);

        sec.append(row);

        return f
    }

    f.addSection = ({name, ...options}: AddSectionOptions): Form<T> => {
        const isSubsection = !!options?.subsectionOf;

        if (isSubsection && !sections[options.subsectionOf]) {
            throw new Error(`Subsection parent with name ${options.subsectionOf} not found`);
        }

        if (sections[name]) {
            throw new Error(`Section with name ${name} already exists`);
        }

        const section = createSectionRow(options);

        sections[name] = section;

        if (isSubsection) {
            sections[options.subsectionOf].append(section);
        } else {
            formElement.append(section);
        }

        return f
    }

    f.onSubmit = (callback: (data: T) => void): Form<T> => {
        submitHandlers.push(callback);

        return f
    }

    f.onChange = (callback: (data: T) => void): Form<T> => {
        changeHandlers.push(callback);

        return f
    }

    f.submit = (): void => {
        const data = options.transformer(new FormData(formElement));

        const validationErrors = options.validator?.(data) ?? [];

        console.log(validationErrors);

        renderValidationErrors(formElement, validationErrors);

        if (validationErrors.length === 0) {
            submitHandlers.forEach(handler => handler(data));
        }
    }

    return f
}

type HiddenField = { type: 'hidden'; value: string; }
type TextField = { type: 'text'; value: string; }
type NumberField = { type: 'number'; value: number; max?: number; min?: number; step?: number; }
type SelectField = { type: 'select'; options: string[]; value: string; }
type CheckboxField = { type: 'checkbox'; value: boolean; }
type TextAreaField = { type: 'textarea'; value: string; }

type FieldWidth = 'narrow' | 'standard' | 'wide' | 'full';

type FieldOptions =
    { label: string; name: string, width?: FieldWidth, mono?: boolean }
    & (HiddenField | TextField | NumberField | SelectField | CheckboxField | TextAreaField);

const createFieldRow = (onChangeHandler: () => void, label: string, ...options: DistributiveOmit<FieldOptions, 'label'>[]): HTMLElement => {
    if (options.every(o => o.type === 'hidden')) {
        // If all the fields are hidden, we don't want a row with all applicable styling.
        return elementFromTemplate(`
            <div class="contents">
                ${options.map(o => `<input type="hidden" name="${options[0].name}" value="${(options[0] as HiddenField).value}">`).join('\n')}
            </div>
        `);
    }

    const row = elementFromTemplate(`
        <div class="sm:grid sm:grid-cols-6 sm:items-start sm:gap-4 sm:py-3">
          <label for="first-name" class="block text-sm font-medium leading-6 text-gray-900 sm:pt-1.5">${label}</label>
          <div class="mt-2 sm:col-span-5 sm:mt-0 flex gap-x-4 js-input">
          </div>
        </div>
    `);

    for (const option of options) {
        if (option.type === 'hidden') {
            return elementFromTemplate(`
                <input type="hidden" name="${option.name}" value="${option.value}">
            `);
        }

        let input: HTMLElement;

        let width = 'sm:max-w-xs';

        switch (option.width) {
            case 'narrow':
                width = 'sm:max-w-44';
                break;
            case 'standard':
                width = 'sm:max-w-xs';
                break;
            case 'wide':
                width = 'sm:max-w-lg';
                break;
            case 'full':
                width = 'sm:max-w-full';
                break;
        }

        switch (option.type) {
            case 'text':
                input = elementFromTemplate(`
                    <input type="text" name="${option.name}" value="${option.value}" class="${width} ${option.mono ? 'font-mono' : ''} block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-blue-600 sm:text-sm sm:leading-6">
                `);
                break;
            case 'number':
                input = elementFromTemplate(`
                    <input type="number" name="${option.name}" class="${width} ${option.mono ? 'font-mono' : ''} block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-blue-600 sm:text-sm sm:leading-6">
                `);
                break;
            case 'select':
                input = elementFromTemplate(`
                    <select name="${option.name}" class="${width} ${option.mono ? 'font-mono' : ''} block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-blue-600 sm:text-sm sm:leading-6">
                        ${option.options.map(opt => `<option value="${opt}" ${opt === option.value ? 'selected' : ''}>${opt}</option>`).join('')}
                    </select>
                `);
                break;
            case 'checkbox':
                input = elementFromTemplate(`
                    <input type="checkbox" name="${option.name}" ${option.value ? 'checked' : ''} class="${width} ${option.mono ? 'font-mono' : ''} block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-blue-600 sm:text-sm sm:leading-6">
                `);
                break;
            case 'textarea':
                input = elementFromTemplate(`
                    <textarea name="${option.name}" class="${width} ${option.mono ? 'font-mono' : ''} block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-blue-600 sm:text-sm sm:leading-6" rows="4">${option.value}</textarea>
                `);
                break;
        }

        input.addEventListener('change', onChangeHandler);

        row.querySelector('.js-input')?.append(input);
    }

    return row
}

const createSectionRow = (options: Omit<AddSectionOptions, 'name'>): HTMLElement => {
    return elementFromTemplate(`
        <div class="pt-8 ${!!options.subsectionOf ? 'contents' : ''}">
            ${createSectionTitle(options.label, options.description)}
        </div>
    `);
}

const createSectionTitle = (label: string | undefined, description: string | undefined): string => {
    return label === undefined
        ? ''
        : `
            <div class="-ml-2 -mt-2 mb-8">
              <h3 class="ml-2 mt-2 text-base font-semibold leading-6 text-gray-900">${label}</h3>
              ${description !== undefined ? `<p class="ml-2 mt-1 truncate text-sm text-gray-500">${description.replace(/\.+$/, '')}.</p>` : ''}
            </div>
        `;
}

const renderValidationErrors = (element: HTMLElement, errors: ValidationError[]): void => {
    element.querySelector('.js-errors')?.remove();

    if (errors.length === 0) {
        return;
    }

    const errorList = elementFromTemplate(`
        <div class="overflow-hidden rounded-lg bg-white shadow js-errors">
          <div class="px-4 py-5 sm:p-6">
            <h3 class="text-sm font-medium text-red-800 mb-4">There were errors with this submission:</h3>
            <table>
                ${errors.map(({ field, message }) => `<tr><td class="text-sm pr-3">${field}</td><td class="text-sm">${message}</td></tr>`).join('')}
            </table>
          </div>
        </div>
    `);

    element.prepend(errorList);
}
