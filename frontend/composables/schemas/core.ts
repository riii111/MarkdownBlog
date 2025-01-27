import type { ValiError, BaseSchema, BaseIssue } from "valibot";

export interface ValidationResult<
  T,
  TSchema extends BaseSchema<unknown, unknown, BaseIssue<unknown>>
> {
  success: boolean;
  data: T | null;
  error: ValiError<TSchema> | null;
}
