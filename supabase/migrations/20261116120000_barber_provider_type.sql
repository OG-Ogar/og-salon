-- Barbers and hair stylists share one account system, but the home feed needs
-- to label their posts with the role each provider actually registered as.
alter table public.barbers
  add column if not exists provider_type text not null default 'barber'
  check (provider_type in ('barber', 'hairstylist'));
