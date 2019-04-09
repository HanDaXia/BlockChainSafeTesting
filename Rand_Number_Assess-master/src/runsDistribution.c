#include <stdio.h>
#include <stdlib.h>
#include <math.h>
#include <string.h>
#include "../include/externs.h"
#include "../include/cephes.h"

//�γ̷ֲ����
void RunsDistribution(int n)
{
	int i = 0;
	int j = 0;
	int num = 0;
	int last = 0; //ǰһλbֵ
	unsigned int k = 0; //ei >= 5�������
	unsigned int *b = NULL; //��Ԫ�����г���Ϊi��1�γ̵���Ŀ
	unsigned int *g = NULL; //��Ԫ�����г���Ϊi��0�γ̵���Ŀ
	double e = 0.0;
	double sum = 0.0;
	double dPValue = 0.0;
	
	for(i = 0;i < n; ++i)
	{
		//
		if(((n-i+3)/pow(2, i+2)) >= 5)
		{
			k = i;
		}
	}
	
	b = (unsigned int*)calloc(k, sizeof(unsigned int));
	if(b == NULL)
	{
		printf("\tRUNSDISTRIBUTION TESTS ABORTED DUE TO ONE OF THE FOLLWING : \n");
    printf("\tInsufficient memory for calloc b. \n");
		return ;
	}

	g = (unsigned int*)calloc(k, sizeof(unsigned int));
	if(g == NULL)
	{
		if(b != NULL) free(b);
		printf("\tRUNSDISTRIBUTION TESTS ABORTED DUE TO ONE OF THE FOLLWING : \n");
    printf("\tInsufficient memory for calloc g. \n");
		return ;
	}
	
	for(i = 1;i <= k; ++i)
	{
		b[i-1] = 0;
		g[i-1] = 0;
	}
	
	last = 0;

		for(j = 0;j < n; ++j)
		{
			if(epsilon[j] == 1)
			{
				//RL_DEBUG("RunsDistribution: pbsTestStream->bits[%d].b:[%d]", j, 1);
				
				if(last != 1)
				{
					if(num >= 1 && num <= k)
					{
						++g[num-1]; //
						
						//RL_DEBUG("RunsDistribution: ++g, g:[%d], i:[%d], j:[%d]", g, i, j);
					}

					//1�γ̿�ʼ
					num = 1;
				}
				else
				{
					++num;
				}

				last = 1; //��¼��һλֵ
			}
			else //pbsTestStream->bits[j].b == 0
			{
				//RL_DEBUG("RunsDistribution: pbsTestStream->bits[%d].b:[%d]", j, 0);
				
				if(last != 0)
				{
					if(num >= 1 && num <= k)
					{
						++b[num-1]; //

						//RL_DEBUG("RunsDistribution: ++b, b:[%d], i:[%d], j:[%d]", b, i, j);
					}
					//0�γ̿�ʼ
					num = 1;
				}
				else
				{
					++num;
				}

				last = 0; //��¼��һλֵ
			}
		}//end for(j = 0;j < pbsTestStream->bitsNumber; ++j)
	
	for(i = 1;i <= k; ++i)
	{
		//RL_DEBUG("RunsDistribution: i : [%d]", i);
		//RL_DEBUG("RunsDistribution: b[%d] : [%d]", i-1, b[i-1]);
		//RL_DEBUG("RunsDistribution: g[%d] : [%d]", i-1, g[i-1]);
		
		e = (double)(n-i+3)/pow(2, i+2);
		
		sum += pow(b[i-1]-e, 2)/e + pow(g[i-1]-e, 2)/e;
	}//end for(i = 1;i <= k; ++i)
	
	dPValue = cephes_igamc(k - 1, sum / (double)2);
	//dPValue = 0;
	
	if(b != NULL) free(b);
	if(g != NULL) free(g);
	
  /*fprintf(stats[TEST_RUNSDISTRIBUTION], "\t\tRUNS DISTRIBUTION TEST\n");
  fprintf(stats[TEST_RUNSDISTRIBUTION], "\t\t------------------------------------------\n");
  fprintf(stats[TEST_RUNSDISTRIBUTION], "\t\tCOMPUTATIONAL INFORMATION:\n");
  fprintf(stats[TEST_RUNSDISTRIBUTION], "\t\t------------------------------------------\n");
  fprintf(stats[TEST_RUNSDISTRIBUTION], "\t\t(a) k                        = %d\n", k);
  fprintf(stats[TEST_RUNSDISTRIBUTION], "\t\t(b) sum                      = %f\n", sum);
  fprintf(stats[TEST_RUNSDISTRIBUTION], "\t\t------------------------------------------\n");

  fprintf(stats[TEST_RUNSDISTRIBUTION], "%s\t\t p_value = %f\n", dPValue < ALPHA ? "FAILURE" : "SUCCESS", dPValue);
  fflush(stats[TEST_RUNSDISTRIBUTION]);
  fprintf(results[TEST_RUNSDISTRIBUTION], "%f\n", dPValue);
  fflush(results[TEST_RUNSDISTRIBUTION]);  */
	if (stats[TEST_RUNSDISTRIBUTION] != 1) {
		stats[TEST_RUNSDISTRIBUTION] = dPValue < ALPHA ? 1 : 0;
	}
  return;
}
