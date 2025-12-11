CREATE TABLE subscription (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "plan" TEXT NOT NULL,
    "referenceId" UUID NOT NULL,
    "stripeCustomerId" TEXT,
    "stripeSubscriptionId" TEXT,
    "status" TEXT NOT NULL,
    "periodStart" TIMESTAMPTZ,
    "periodEnd" TIMESTAMPTZ,
    "cancelAtPeriodEnd" BOOLEAN,
    "seats" INTEGER,
    "trialStart" TIMESTAMPTZ,
    "trialEnd" TIMESTAMPTZ,
);
