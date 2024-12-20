import { Validator } from './input';

export type Validation = {
    maxLength?: number;
    minLength?: number;
    regex?: RegExp;
    required?: true,
    // Async validations
    remote?: 'jsonata' | 'template',
}

export const validator = (validation: Validation, onValidate: (valid: boolean) => void): Validator => {
    // Add AbortControllers for other types of async validations.
    let jsonataAbortController: AbortController | undefined;

    const err = (message: string | null) => {
        onValidate(message === null);
        return message;
    }

    return async (value: string): Promise<string | null> => {
        const cleanedValue = value.trim();

        // Keep validations in order of efficiency. No point in calling a remove endpoint, when we can already fail
        // validation by just looking at the string.

        if (!validation?.required && cleanedValue === '') {
            return err(null);
        }

        if (validation?.required && cleanedValue === '') {
            return err('This field is required');
        }

        if (validation?.minLength && cleanedValue.length < validation.minLength) {
            return err(`This field must be at least ${ validation.minLength } characters`);
        }

        if (validation?.maxLength && cleanedValue.length > validation.maxLength) {
            return err(`This field must be less than ${ validation.maxLength } characters`);
        }

        if (validation?.regex && !new RegExp(validation.regex).test(cleanedValue)) {
            return err(`This field should match the pattern ${ validation.regex }`);
        }

        if (validation?.remote && cleanedValue !== '') {
            if (jsonataAbortController) {
                jsonataAbortController.abort();
            }

            jsonataAbortController = new AbortController();

            const resp = await fetch('/api/v1/validate', {
                signal: jsonataAbortController.signal,
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ type: validation.remote, subject: cleanedValue }),
            });

            if (!resp.ok) {
                return err('Could not validate this input as jsonata.');
            }

            const body = await resp.json();

            return err(( body && body.error ) || null);
        }

        return err(null);
    }
}
