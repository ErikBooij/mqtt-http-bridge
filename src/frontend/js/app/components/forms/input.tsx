import React, { FormEvent, FormEventHandler, useState } from 'react'

type CommonProps = {
    onChange?: (value: string) => void,
    validator?: Validator,
    value?: string | undefined,
}

export type TextFieldProps = CommonProps & {
    type?: 'text' | 'password' | 'email' | 'number',
}

export type Validator = (value: string) => Promise<string | null>;

export const TextField = ({ onChange, validator, value }: TextFieldProps) => {
    const [ error, setError ] = useState<string | null>(null);

    return (
        <>
            <input type="text"
                   name="name"
                   className="block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-blue-600 sm:text-sm sm:leading-6"
                   defaultValue={ value }
                   onInput={ configureValidation(validator, setError, onChange) }
            />
            { error !== null && <div className="js-error text-red-600 text-sm mt-2">{ error }</div> }
        </>
    )
}

type TextAreaProps = CommonProps & {
    rows?: number
}

export const TextArea = ({ onChange, rows, validator, value }: TextAreaProps) => {
    const [ error, setError ] = useState<string | null>(null);

    return (
        <>
            <textarea
                className="font-mono block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-blue-600 sm:text-sm sm:leading-6"
                onInput={ configureValidation(validator, setError, onChange) }
                rows={ rows || 4 }
                defaultValue={ value }
            />
            { error !== null && <div className="js-error text-red-600 text-sm mt-2">{ error }</div> }
        </>
    )
}

type SelectProps = CommonProps & {
    options: Record<string, string>
}

export const Select = ({ onChange, options, validator, value }: SelectProps) => {
    const [ error, setError ] = useState<string | null>(null);

    return <>
        <select
            className="sm:max-w-44 block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-blue-600 sm:text-sm sm:leading-6"
            defaultValue={ value }
            name="method"
            onInput={ configureValidation(validator, setError, onChange) }
        >
            { Object.entries(options).map(([ value, label ]) => (
                <option key={ value } value={ value }>{ label }</option>
            )) }
        </select>
        { error !== null && <div className="js-error text-red-600 text-sm mt-2">{ error }</div> }
    </>
}

const configureValidation = <T extends HTMLElement>(validator: Validator | undefined, setError: (error: string | null) => void, onChange: ( (e: string) => unknown ) | undefined): FormEventHandler<T> => {
    return async (e: FormEvent<T>): Promise<void> => {
        const cleanedValue = ( e.target as HTMLInputElement )?.value.trim();

        onChange?.(cleanedValue);

        if (validator) {
            const error = await validator(cleanedValue);

            setError(error);
        }
    }
}
