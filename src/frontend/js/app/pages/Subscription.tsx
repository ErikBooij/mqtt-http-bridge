import React, { PropsWithChildren, useContext, useEffect, useState } from 'react'
import { LayoutContext } from '../components/Layout';
import { useNavigate, useParams } from 'react-router';
import { useCreateSubscription, useFetchSubscription, useUpdateSubscription } from '../../rq/subscription';
import { PageTitle } from '../components/PageTitle';
import { Select, TextArea, TextField, Validator } from '../components/forms/input';
import { SubscriptionWithoutID as APISubscription } from '../../api/subscriptions';
import { Validation, validator } from '../components/forms/validation';

type Props = {
    mode: 'new' | 'copy' | 'edit';
}

export const Subscription = ({ mode }: Props) => {
    const { setCurrentPage } = useContext(LayoutContext);

    useEffect(() => {
        setCurrentPage('subscriptions');
    })

    const [ name, setName ] = useState('');
    const [ topic, setTopic ] = useState('');
    const [ extract, setExtract ] = useState<Record<string, string>>({});
    const [ filter, setFilter ] = useState('');
    const [ method, setMethod ] = useState('POST');
    const [ url, setUrl ] = useState('');
    const [ headers, setHeaders ] = useState<Record<string, string>>({});
    const [ body, setBody ] = useState('');
    const [ validations, setValidations ] = useState<Record<string, boolean>>({});

    const navigate = useNavigate();

    const afterSave = () => {
        navigate('/');
    }

    const createSubscription = useCreateSubscription({ onSuccess: afterSave });
    const updateSubscription = useUpdateSubscription({ onSuccess: afterSave });

    const { id } = useParams();

    const shouldFetch = mode === 'copy' || mode === 'edit' && !!id;

    const { isFetching, error, data: subscription } = useFetchSubscription(id || '', shouldFetch);

    useEffect(() => {
        if (subscription) {
            setName(subscription.name);
            setTopic(subscription.topic);
            setExtract(subscription.extract || {});
            setFilter(subscription.filter || '');
            setMethod(subscription.method);
            setUrl(subscription.url);
            setHeaders(subscription.headers || {});
            setBody(subscription.body || '');
        }
    }, [ subscription ]);

    if (isFetching) {
        return <div>Loading...</div>;
    }

    if (error) {
        return <div>Error: { error.message }</div>;
    }

    const buildSubscriptionObject = (): APISubscription => ( {
        name,
        topic,
        extract,
        filter,
        method: method as APISubscription['method'],
        url,
        headers,
        body,
    } )

    const saveSubscription = async () => {
        if (mode === 'new' || mode === 'copy') {
            createSubscription.mutate({ subscription: buildSubscriptionObject() });
        } else {
            updateSubscription.mutate({ id: id!, subscription: buildSubscriptionObject() });
        }
    }

    const registerValidation = (field: string, validation: Validation): Validator => {
        return validator(validation, (valid) => {
            setValidations(v => ( { ...v, [ field ]: valid } ));
        })
    }

    return (
        <div>
            <PageTitle action={ {
                title: 'Save',
                onClick: saveSubscription,
                disabled: Object.values(validations).some(v => !v)
            } }>Subscription</PageTitle>
            <div className="sm:grid sm:grid-cols-12">
                <div className="sm:col-span-8">
                    <form
                        className="space-y-8 border-b pb-12 sm:space-y-8 sm:divide-gray-900/10 sm:divide-y-2 sm:pb-16">
                        <Section>
                            <FieldRow label="Name">
                                <TextField value={ name } onChange={ setName }
                                           validator={ registerValidation('name', { required: true }) }/>
                            </FieldRow>
                        </Section>
                        <Section>
                            <SectionHeading label="Input">
                                Configure which messages to act on, and how to extract data from them.
                            </SectionHeading>
                            <FieldRow label="Topic">
                                <TextField value={ topic } onChange={ setTopic }
                                           validator={ registerValidation('topic', { required: true }) }
                                />
                            </FieldRow>
                            <FieldRow label="Extract">
                                <ExpandingKeyValue value={ subscription?.extract } onChange={ setExtract }
                                                   keyValidation={ (i: number) => registerValidation(`extract_key_${ i }`, { regex: /^[a-zA-Z0-9_\-.]+$/ }) }
                                                   valueValidation={ (i: number) => registerValidation(`extract_value_${ i }`, { remote: 'jsonata' }) }
                                />
                            </FieldRow>
                            <FieldRow label="Filter">
                                <TextField value={ filter } onChange={ setFilter }
                                           validator={ registerValidation('filter', { remote: 'jsonata' }) }
                                />
                            </FieldRow>
                        </Section>
                        <Section>
                            <SectionHeading label="Output">
                                The HTTP request to send when a message is received.
                            </SectionHeading>
                            <FieldRow label="Method">
                                <Select value={ method } onChange={ setMethod } options={ {
                                    'GET': 'GET',
                                    'HEAD': 'HEAD',
                                    'POST': 'POST',
                                    'PUT': 'PUT',
                                    'PATCH': 'PATCH',
                                    'DELETE': 'DELETE',
                                } }/>
                            </FieldRow>
                            <FieldRow label="URL">
                                <TextField value={ url } onChange={ setUrl }
                                           validator={ registerValidation('url', { remote: 'template' }) }
                                />
                            </FieldRow>
                            <FieldRow label="Headers">
                                <ExpandingKeyValue value={ subscription?.headers } onChange={ setHeaders }
                                                   keyValidation={ (i: number) => registerValidation(`header_key_${ i }`, { regex: /^[a-zA-Z0-9_\-.]+$/ }) }
                                                   valueValidation={ (i: number) => registerValidation(`header_value_${ i }`, { remote: 'template' }) }
                                />
                            </FieldRow>
                            <FieldRow label="Body">
                                <TextArea value={ body } onChange={ setBody }
                                          validator={ registerValidation('body', { required: true }) }
                                />
                            </FieldRow>
                        </Section>
                    </form>
                </div>
            </div>
        </div>
    );
}

const ExpandingKeyValue = ({ onChange, value, keyValidation, valueValidation }: {
    onChange?: (e: Record<string, string>) => void,
    value: Record<string, string> | undefined,
    keyValidation?: (i: number) => Validator,
    valueValidation?: (i: number) => Validator,
}) => {
    const [ rows, setRows ] = useState<[ string, string ][]>([ ...Object.entries(value || {}), [ '', '' ] ]);

    const reportLatestState = () => {
        onChange?.(Object.fromEntries(rows.filter(([ key, value ]) => key && value)));
        console.log({ rows });
    }

    const onChangeKey = (i: number) => (e: string) => {
        setRows(r => {
            const newValue = [ ...r ];

            newValue[ i ] = [ e, newValue[ i ][ 1 ] || '' ];

            if (i === newValue.length - 1) {
                newValue.push([ '', '' ]);
            }

            return newValue;
        })

        reportLatestState();
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

        reportLatestState();
    }

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

const FieldRow = ({ label, children }: PropsWithChildren<{ label: string }>) => {
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

const Section = ({ children }: PropsWithChildren) => {
    return <div className="pt-8">
        { children }
    </div>
}

const SectionHeading = ({ children, label }: PropsWithChildren<{ label: string }>) => {
    return <div className="-ml-2 -mt-2 mb-8">
        <h3 className="ml-2 mt-2 text-base font-semibold leading-6 text-gray-900">{ label }</h3>
        <p className="ml-2 mt-1 truncate text-sm text-gray-500">
            { children }
        </p>
    </div>
}
