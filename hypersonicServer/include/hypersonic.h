#ifndef HYPERSONIC
#define HYPERSONIC


typedef struct
{
    int channel;
    int attenuation;
} hypersonic_config;

void initialise_hypersonic();
int get_distance(hypersonic_config *config);

void sample_hypersonic_task(void *params);





#endif