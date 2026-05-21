CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE events (
                        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                        title VARCHAR(255) NOT NULL,
                        date TIMESTAMP NOT NULL,
                        location VARCHAR(255) NOT NULL
);

CREATE TABLE event_registrations (
                                     id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                     event_id UUID REFERENCES events(id) ON DELETE CASCADE,
                                     student_id UUID NOT NULL,
                                     UNIQUE(event_id, student_id)
);

CREATE TABLE coworking_bookings (
                                    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                    space_id UUID NOT NULL,
                                    student_id UUID NOT NULL,
                                    start_time TIMESTAMP NOT NULL,
                                    end_time TIMESTAMP NOT NULL
);