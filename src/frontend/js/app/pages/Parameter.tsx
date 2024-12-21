import React, { useEffect, useState } from 'react'
import { PageTitle } from '../components/PageTitle';
import { useNavigate, useParams } from 'react-router';
import { useDeleteGlobalParameter, useListGlobalParameters, useSetGlobalParameter } from '../../rq/parameter';
import { FieldRow, Section } from '../components/FormSections';
import { TextField } from '../components/forms/input';

export const Parameter = () => {
    const navigate = useNavigate();

    const globalParametersQuery = useListGlobalParameters();
    const saveParameter = useSetGlobalParameter({});
    const deleteParameter = useDeleteGlobalParameter({});

    const [ key, setKey ] = useState('');
    const [ value, setValue ] = useState('');

    const { key: initialKey } = useParams();

    if (initialKey) {
        useEffect(() => {
            setKey(initialKey);
            setValue(globalParametersQuery.data?.[ initialKey ] || '');
        }, [ initialKey, globalParametersQuery.data ]);
    }

    if (globalParametersQuery.isPending) {
        return <div>Loading...</div>;
    }

    if (globalParametersQuery.error) {
        return <div>Error: { globalParametersQuery.error.message }</div>;
    }

    const saveAction = () => {
        const onSuccess = () => navigate('/parameters');
        const doSave = (after: () => void) => saveParameter.mutate({ key, value }, { onSuccess: after });

        if (initialKey && initialKey !== key) {
            deleteParameter.mutate({ key: initialKey }, { onSuccess: () => doSave(onSuccess) });
            return
        }

        doSave(onSuccess);
    };

    const deleteAction = () => {
        if (initialKey) {
            deleteParameter.mutate({ key: initialKey }, { onSuccess: () => navigate('/parameters') });
        }
    }

    return (
        <div>
            <PageTitle
                currentPage="parameters"
                action={ { title: 'Save', onClick: saveAction } }
                secondaryAction={ initialKey ? { title: 'Delete', onClick: deleteAction } : undefined }
            >
                Global Parameters
            </PageTitle>
            <form
                className="space-y-8 border-b pb-12 sm:space-y-8 sm:divide-gray-900/10 sm:divide-y-2 sm:pb-16">
                <Section>
                    <FieldRow label="Key">
                        <TextField value={ key } onChange={ setKey }/>
                    </FieldRow>
                    <FieldRow label="Value">
                        <TextField value={ value } onChange={ setValue }/>
                    </FieldRow>
                </Section>
            </form>
        </div>
    );
}
