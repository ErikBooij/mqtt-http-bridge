import {
    apiErrorFromResponse,
    AsyncMaybeAPIError,
    responseParseError,
} from "./common";
import {z} from "zod";

const subscriptionSchema = z.object({
    id: z.string().uuid(),
    name: z.string().min(1),
    topic: z.string().min(1),
    extract: z.record(z.string(), z.string()).optional(),
    filter: z.string().optional(),
    method: z.enum(['GET', 'POST', 'PATCH', 'PUT', 'DELETE', 'HEAD', 'OPTIONS']),
    url: z.string(),
    headers: z.record(z.string(), z.string()).optional(),
    body: z.string().optional(),
    templateId: z.string().uuid().optional(),
    templateParameters: z.map(z.string(), z.any()).optional(),
}).strict();

const subscriptionResponseSchema = z.object({
    subscription: subscriptionSchema,
});

const subscriptionsResponseSchema = z.object({
    subscriptions: z.array(subscriptionSchema),
});

export type Subscription = z.infer<typeof subscriptionSchema>;
export type SubscriptionWithoutID = Omit<Subscription, 'id'>;
export type SubscriptionResponse = z.infer<typeof subscriptionResponseSchema>; // Single subscription
export type SubscriptionsResponse = z.infer<typeof subscriptionsResponseSchema>; // Multiple subscriptions

export const fetchSubscription = async (id: string): AsyncMaybeAPIError<Subscription> => {
    const response = await fetch(`/api/v1/subscriptions/${id}`);

    if (response.status !== 200) {
        return [null, await apiErrorFromResponse(response)];
    }

    const parsedResponse = subscriptionResponseSchema.safeParse(await response.json());

    if (!parsedResponse.success) {
        return [null, responseParseError(parsedResponse.error)];
    }

    return [parsedResponse.data.subscription, null];
}

export const fetchSubscriptions = async (): AsyncMaybeAPIError<Subscription[]> => {
    const response = await fetch('/api/v1/subscriptions');

    if (response.status !== 200) {
        return [null, await apiErrorFromResponse(response)];
    }

    const parsedResponse = subscriptionsResponseSchema.safeParse(await response.json());

    if (!parsedResponse.success) {
        return [null, responseParseError(parsedResponse.error)];
    }

    return [parsedResponse.data.subscriptions, null];
}


export const createSubscription = async (subscription: SubscriptionWithoutID): AsyncMaybeAPIError<Subscription> => {
    const response = await fetch('/api/v1/subscriptions', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(subscription),
    });

    if (response.status !== 201) {
        return [null, await apiErrorFromResponse(response)];
    }

    const parsedResponse = subscriptionResponseSchema.safeParse(await response.json());

    if (!parsedResponse.success) {
        return [null, responseParseError(parsedResponse.error)];
    }

    return [parsedResponse.data.subscription, null];
}

export const updateSubscription = async ({ id, ...subscription }: Subscription): AsyncMaybeAPIError<Subscription> => {
    const response = await fetch(`/api/v1/subscriptions/${id}`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(subscription),
    });

    if (response.status !== 201) {
        return [null, await apiErrorFromResponse(response)];
    }

    const parsedResponse = subscriptionResponseSchema.safeParse(await response.json());

    if (!parsedResponse.success) {
        return [null, responseParseError(parsedResponse.error)];
    }

    return [parsedResponse.data.subscription, null];
}

export const deleteSubscription = async (id: string): AsyncMaybeAPIError<void> => {
    const response = await fetch(`/api/v1/subscriptions/${id}`, {
        method: 'DELETE',
    });

    if (response.status !== 200) {
        return [null, await apiErrorFromResponse(response)];
    }

    return [void 0, null];
}
