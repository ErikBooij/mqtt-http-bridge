import React, { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router';
import {
    useCreateSubscription,
    useDeleteSubscription,
    useFetchSubscription,
    useUpdateSubscription
} from '../../rq/subscription';
import { PageTitle } from '../components/PageTitle';
import { Select, TextArea, TextField, Validator } from '../components/forms/input';
import { SubscriptionWithoutID as APISubscription } from '../../api/subscriptions';
import { Validation, validator } from '../components/forms/validation';
import { useListGlobalParameters } from '../../rq/parameter';
import { CopyIcon } from '../components/Icons';
import { ExpandingKeyValue, FieldRow, Section, SectionHeading } from '../components/FormSections';

type Props = {
    mode: 'new' | 'copy' | 'edit';
}

export const Subscription = ({ mode }: Props) => {
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

    const globalParameters = useListGlobalParameters();

    const [ parameters, setParameters ] = useState<{ extract: string[], global: string[], meta: string[] }>({
        extract: [],
        global: [],
        meta: [
            'topic',
            'client',
            'payload',
        ],
    });

    useEffect(() => {
        if (globalParameters.data) {
            setParameters(p => ( {
                ...p,
                global: Object.keys(globalParameters.data),
            } ));
        }

        if (extract) {
            setParameters(p => ( {
                ...p,
                extract: Object.keys(extract),
            } ));
        }
    }, [ globalParameters.data, extract ]);

    const createSubscription = useCreateSubscription({ onSuccess: afterSave });
    const updateSubscription = useUpdateSubscription({ onSuccess: afterSave });
    const deleteSubscription = useDeleteSubscription();

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
            <PageTitle
                action={ {
                    title: 'Save',
                    onClick: saveSubscription,
                    disabled: Object.values(validations).some(v => !v)
                } }
                secondaryAction={ mode === 'edit' ? {
                    title: 'Delete',
                    onClick: () => {
                        if (confirm('Are you sure you want to delete this subscription?')) {
                            deleteSubscription.mutate({ id: id! }, { onSuccess: afterSave });
                        }
                    }
                } : undefined}
                currentPage="subscriptions"
            >Subscription</PageTitle>
            <div className="flex flex-col-reverse sm:grid sm:grid-cols-12 sm:gap-x-8">
                <div className="sm:col-span-8">
                    <form
                        className="space-y-8 border-b pb-12 sm:space-y-8 sm:divide-gray-900/10 sm:divide-y-2 sm:pb-16">
                        <Section>
                            <FieldRow label="Name">
                                <TextField value={ name } onChange={ setName }
                                           validator={ registerValidation('name', {
                                               required: true,
                                               regex: /^[0-9A-Za-z]([0-9A-Za-z.\-_|\\/ ]*[0-9A-Za-z])?$/
                                           }) }/>
                            </FieldRow>
                        </Section>
                        <Section>
                            <SectionHeading label="Input">
                                Configure which messages to act on, and how to extract data from them.
                            </SectionHeading>
                            <FieldRow label="Topic">
                                <TextField value={ topic } onChange={ setTopic }
                                           validator={ registerValidation('topic', {
                                               required: true,
                                               regex: /^(\/?(([^#+\\/]*|\+)(\/([^#+\\/]*|\+))*)?\/?#?)$/
                                           }) }
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
                <div className="sm:col-span-4 pt-8">
                    <div className="bg-blue-200/50 p-4 rounded-md">
                        <h3 className="text-sm font-medium">Parameters</h3>
                        {
                            Object.entries(parameters)
                                .sort(([ sectionA ], [ sectionB ]) => sectionA.localeCompare(sectionB))
                                .map(([ section, keys ]) => (
                                    <div key={ section } className="mt-4">
                                        <ul className="overflow-hidden rounded-md bg-white px-4 py-4 shadow">
                                            {
                                                keys.sort(([ keyA ], [ keyB ]) => keyA.localeCompare(keyB)).map(key => (
                                                    <li key={ key }
                                                        className="whitespace-nowrap flex items-center gap-x-2 py-1 cursor-pointer hover:bg-gray-100 -mx-2 px-2 rounded-md"
                                                        onClick={ () => navigator.clipboard.writeText(`{{ .${ section }.${ key } }}`) }
                                                    >
                                                        { <CopyIcon/> }
                                                        <span className="font-mono text-sm">
                                                        { section }
                                                            .
                                                        <span className="font-bold">{ key }
                                                    </span>
                                            </span></li>
                                                ))
                                            }
                                        </ul>
                                    </div>
                                ))
                        }
                    </div>
                </div>
            </div>
        </div>
    );
}
