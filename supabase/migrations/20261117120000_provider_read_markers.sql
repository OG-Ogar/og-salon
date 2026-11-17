-- Remember what each provider has already opened so the home shortcuts can
-- show new bookings and unread messages instead of lifetime totals.
create table if not exists degosalon_internal.provider_read_markers (
    barber_id bigint primary key references public.barbers(id) on delete cascade,
    booking_id bigint not null default 0,
    updated_at timestamptz not null default now()
);

create table if not exists degosalon_internal.chat_reads (
    barber_id bigint not null references public.barbers(id) on delete cascade,
    customer_id bigint not null references public.customers(id) on delete cascade,
    last_message_id bigint not null default 0,
    updated_at timestamptz not null default now(),
    primary key (barber_id, customer_id)
);

create index if not exists chat_reads_owner_idx
    on degosalon_internal.chat_reads(barber_id);

-- Existing activity is the starting point, not a backlog. Only bookings,
-- orders and messages that arrive after this migration raise a shortcut.
insert into degosalon_internal.provider_read_markers(barber_id, booking_id)
select b.id, coalesce((select max(x.id) from public.bookings x where x.barber_id = b.id), 0)
from public.barbers b
on conflict (barber_id) do nothing;

insert into degosalon_internal.chat_reads(barber_id, customer_id, last_message_id)
select barber_id, customer_id, max(id)
from public.messages
group by barber_id, customer_id
on conflict (barber_id, customer_id) do nothing;

revoke all on degosalon_internal.provider_read_markers from public;
revoke all on degosalon_internal.chat_reads from public;
alter table degosalon_internal.provider_read_markers enable row level security;
alter table degosalon_internal.chat_reads enable row level security;
