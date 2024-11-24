export const validateExtract = (abortController: AbortController, extract: string) => {
    return performValidationRequest(abortController, 'extract', extract);
}

export const validateFilter = (abortController: AbortController, filter: string) => {
    return performValidationRequest(abortController, 'filter', filter);
}

const performValidationRequest = async (abortController: AbortController, type: 'extract' | 'filter', subject: string): Promise<string|null> => {
    const resp = await fetch('/api/v1/validate', {
        signal: abortController.signal,
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ type, subject }),
    });

    if (!resp.ok) {
        return 'Could not validate this input.';
    }

    const body = await resp.json();

    return (body && body.error) || null;
}
