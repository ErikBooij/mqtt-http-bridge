import {form} from "./form";
import {navigateToGlobalParametersOverview} from "../shared/navigation";
import {elementFromTemplate} from "../dom/template";
import {copyIcon} from "../shared/icons";
import {setGlobalParameter} from "../api/global-parameters";

interface GlobalParameterForm {
    submit(): void;
}

type GlobalParameterFormOptions = {
    key?: string;
    selector: string;
    value: string | undefined;
}

type GlobalParameter = {
    key: string;
    value: string;
}

export const globalParameterForm = (
    {key, selector, value}: GlobalParameterFormOptions
): GlobalParameterForm => {
    const element = document.querySelector(selector);

    if (!element) {
        throw new Error(`Element with selector ${selector} not found`);
    }

    const subForm = form<GlobalParameter>({
        selector,
        transformer: (data: FormData): GlobalParameter => {
            return {
                key: data.get('key') as string,
                value: data.get('value') as string,
            };
        },
        validator: (data: GlobalParameter) => {
            const errors = [];

            if (data.key.trim() === '') {
                errors.push({field: 'Key', message: 'Value is required'});
            }

            if (data.value.trim() === '') {
                errors.push({field: 'Value', message: 'Value is required'});
            }

            return errors;
        }
    });

    subForm.onSubmit(async ({ key, value}: GlobalParameter) => {
        console.log('Setting global parameter', { key, value });

        const [_, error] = await setGlobalParameter(key, value);

        if (error === null) {
            navigateToGlobalParametersOverview();
        }
    });

    subForm
        .addSection({name: 'identifier'})
        .addField({section: 'identifier', label: 'Key', name: 'key', type: 'text', value: key ?? ''})
        .addField({section: 'identifier', label: 'Value', name: 'value', type: 'text', value: value ?? ''});

    return {
        submit: () => {
            subForm.submit();
        }
    };
}
