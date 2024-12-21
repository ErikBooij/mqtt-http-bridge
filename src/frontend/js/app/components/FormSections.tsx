import { TextField, Validator } from './forms/input';
import React, { PropsWithChildren, useEffect, useState } from 'react';

export const ExpandingKeyValue = ({ onChange, value, keyValidation, valueValidation }: {
    onChange?: (e: Record<string, string>) => void,
    value: Record<string, string> | undefined,
    keyValidation?: (i: number) => Validator,
    valueValidation?: (i: number) => Validator,
}) => {
    const [ rows, setRows ] = useState<[ string, string ][]>([ ...Object.entries(value || {}), [ '', '' ] ]);

    const onChangeKey = (i: number) => (e: string) => {
        setRows(r => {
            const newValue = [ ...r ];

            newValue[ i ] = [ e, newValue[ i ][ 1 ] || '' ];

            if (i === newValue.length - 1) {
                newValue.push([ '', '' ]);
            }

            return newValue;
        })
    }

    const onChangeValue = (i: number) => (e: string) => {
        setRows(r => {
            const newValue = [ ...r ];

            newValue[ i ] = [ newValue[ i ][ 0 ] || '', e ];

            if (i === newValue.length - 1) {
                newValue.push([ '', '' ]);
            }

            return newValue;
        })
    }

    useEffect(() => {
        onChange?.(Object.fromEntries(rows.filter(([ key, value ]) => key && value)));
    }, [ rows ]);

    return <div className="flex gap-y-2 flex-col">
        { rows.map(([ key, value ], i) => (
            <div className="w-full block" key={ i }>
                <div className="mt-2 sm:col-span-5 sm:mt-0 flex gap-x-4 js-input">
                    <div className="w-1/4 flex-grow flex-shrink">
                        <TextField value={ key } onChange={ onChangeKey(i) } validator={ keyValidation?.(i) }/>
                    </div>
                    <div className="w-3/4 flex-grow flex-shrink">
                        <TextField value={ value } onChange={ onChangeValue(i) } validator={ valueValidation?.(i) }/>
                    </div>
                </div>
            </div>
        )) }
    </div>
}

export const FieldRow = ({ label, children }: PropsWithChildren<{ label: string }>) => {
    return (
        <div className="sm:grid sm:grid-cols-6 sm:items-start sm:gap-4 sm:py-3">
            <label htmlFor="first-name"
                   className="block text-sm font-medium leading-6 text-gray-900 sm:pt-1.5 sm:mt-0 mt-6"
            >
                { label }
            </label>
            <div className="mt-2 sm:col-span-5 sm:mt-0">
                { children }
            </div>
        </div>
    )
}

export const Section = ({ children }: PropsWithChildren) => {
    return <div className="pt-8">
        { children }
    </div>
}

export const SectionHeading = ({ children, label }: PropsWithChildren<{ label: string }>) => {
    return <div className="-ml-2 -mt-2 mb-8">
        <h3 className="ml-2 mt-2 text-base font-semibold leading-6 text-gray-900">{ label }</h3>
        <p className="ml-2 mt-1 truncate text-sm text-gray-500">
            { children }
        </p>
    </div>
}
