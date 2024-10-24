// Copy-pasted from https://stackoverflow.com/questions/74857965/type-narrowing-not-working-when-using-omit#comment132108966_74857965
// Reason is that Omit does not distribute over unions. You can use DistributiveOmit to achieve this:

export type DistributiveOmit<T, K extends keyof any> = T extends any
    ? Omit<T, K>
    : never;
